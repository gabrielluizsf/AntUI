// Package sdk finds the Android SDK, the NDK and a JDK on the machine doing
// the building, and says where each tool inside them is.
//
// It is host-side: it runs on Linux, macOS and Windows, needs no cgo and no
// device, and knows nothing about a running app. Everything that builds an
// APK goes through here so that "where is aapt2" is answered once.
package sdk

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Toolchain is everything found. [Find] always looks for all of it — a part
// that is simply not installed is zero here, and [Require] decides only
// which of those absences is an error.
type Toolchain struct {
	Root       string // the SDK root
	BuildTools BuildTools
	Platform   Platform
	NDK        NDK
	JDK        JDK
	// Bundletool is the jar that takes an App Bundle apart again, or empty
	// when it is not on this machine. It is **not part of the SDK**: it is a
	// single jar from Google's own releases, and nothing needs it to *build*
	// a bundle — only to check one.
	Bundletool string
}

// BuildTools is one entry under the SDK's build-tools directory: aapt2, d8,
// zipalign and apksigner live here.
type BuildTools struct {
	Dir     string
	Version Version
}

// Platform is one entry under the SDK's platforms directory. It matters for
// two things only: the android.jar the Java shim compiles against, and the
// API level an app says it targets.
type Platform struct {
	Dir  string
	Name string // "android-37.2", as the directory is named
	API  int    // 37
	// Minor is the second number of a platform like android-36.1, which
	// Android started shipping when it began releasing an SDK mid-year. It
	// is 0 on a platform named by a single number.
	Minor int
	// Ext is the extension level of an android-34-ext11, and 0 otherwise.
	// These are not chosen automatically: they are a strictly older platform
	// with newer APIs bolted on, which is not what a build wants by default.
	Ext int
}

// NDK is the native toolchain: the clang that compiles Go's cgo for a phone.
type NDK struct {
	Dir     string
	Version Version
	// Host is the prebuilt directory the toolchain for this machine is in —
	// "linux-x86_64" and the like. It is read off the disk rather than
	// assumed, because which host builds the NDK ships has changed.
	Host string
}

// JDK is a Java compiler, needed only to rebuild the Java shim. Building an
// app does not need one: the shim ships already compiled.
type JDK struct {
	Dir   string
	Javac string
}

// Require says what a caller cannot do without. The zero value asks for a
// usable SDK with build tools and a platform, which is what building an APK
// takes.
type Require struct {
	NDK bool // the app has native code, which every Go app does
	JDK bool // only for rebuilding the Java shim

	// MinBuildTools and MinPlatform are floors, not preferences: the newest
	// found is always the one chosen, and these only decide whether it is
	// good enough.
	MinBuildTools Version
	MinPlatform   int
}

// ErrNoSDK is what [Find] reports when it cannot even locate the SDK, as
// distinct from finding one with pieces missing.
var ErrNoSDK = errors.New("sdk: no Android SDK found")

// Find locates everything Require asks for. The error it returns on a
// partial install is a [*Missing], which names the packages to install and
// the sdkmanager line that installs them.
func Find(req Require) (*Toolchain, error) {
	root := findRoot()
	if root == "" {
		return nil, ErrNoSDK
	}
	t := &Toolchain{Root: root}
	miss := &Missing{Root: root}

	if bt, ok := findBuildTools(root); !ok {
		miss.add("build-tools", "build-tools;37.0.0")
	} else if !bt.Version.AtLeast(req.MinBuildTools) {
		miss.old("build-tools", bt.Version.String(), req.MinBuildTools.String())
		t.BuildTools = bt
	} else {
		t.BuildTools = bt
	}

	if p, ok := findPlatform(root); !ok {
		miss.add("platforms", "platforms;android-37.2")
	} else if p.API < req.MinPlatform {
		miss.old("platforms", p.Name, "android-"+strconv.Itoa(req.MinPlatform))
		t.Platform = p
	} else {
		t.Platform = p
	}

	// Both of these are looked for whether or not they were asked for: a
	// caller that says it does not need a JDK still signs with one if there
	// is one, and finding out later that it was there all along is a worse
	// failure than the search costs.
	if n, ok := findNDK(root); ok {
		t.NDK = n
	} else if req.NDK {
		miss.add("ndk", "ndk;30.0.16138531")
	}
	if j, ok := findJDK(); ok {
		t.JDK = j
	} else if req.JDK {
		miss.add("a JDK", "")
	}
	t.Bundletool = findBundletool(root)
	if len(miss.Problems) > 0 {
		return t, miss
	}
	return t, nil
}

// findRoot looks where the environment says, then where each system puts it
// by default. ANDROID_HOME is the current name and ANDROID_SDK_ROOT the one
// before it; both are still set by real installs, so both are read.
func findRoot() string {
	for _, env := range []string{"ANDROID_HOME", "ANDROID_SDK_ROOT"} {
		if dir := os.Getenv(env); dir != "" && isSDK(dir) {
			return dir
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	var candidates []string
	switch runtime.GOOS {
	case "darwin":
		candidates = []string{filepath.Join(home, "Library", "Android", "sdk")}
	case "windows":
		candidates = []string{filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk")}
	default:
		candidates = []string{
			filepath.Join(home, "Android", "Sdk"),
			filepath.Join(home, "android-sdk"),
			"/opt/android-sdk",
			"/usr/lib/android-sdk",
		}
	}
	for _, dir := range candidates {
		if isSDK(dir) {
			return dir
		}
	}
	return ""
}

// isSDK reports whether a directory looks like an SDK root rather than
// merely existing — an empty ANDROID_HOME pointing at a directory that was
// never populated is a common way to waste an afternoon.
func isSDK(dir string) bool {
	if dir == "" {
		return false
	}
	for _, sub := range []string{"platform-tools", "platforms", "build-tools", "cmdline-tools", "tools", "ndk"} {
		if st, err := os.Stat(filepath.Join(dir, sub)); err == nil && st.IsDir() {
			return true
		}
	}
	return false
}

func findBuildTools(root string) (BuildTools, bool) {
	var found []BuildTools
	for _, name := range subdirs(filepath.Join(root, "build-tools")) {
		dir := filepath.Join(root, "build-tools", name)
		// A build-tools directory that lost its aapt2 is a half-finished
		// download, and picking it would fail much later and less clearly.
		if _, err := os.Stat(filepath.Join(dir, exe("aapt2"))); err != nil {
			continue
		}
		found = append(found, BuildTools{Dir: dir, Version: ParseVersion(name)})
	}
	return pick(found, func(b BuildTools) Version { return b.Version })
}

func findPlatform(root string) (Platform, bool) {
	var found []Platform
	for _, name := range subdirs(filepath.Join(root, "platforms")) {
		p, ok := parsePlatform(name)
		if !ok {
			continue
		}
		p.Dir = filepath.Join(root, "platforms", name)
		if _, err := os.Stat(p.Jar()); err != nil {
			continue
		}
		// An -ext platform is an older API level with later additions, so it
		// is never the right automatic choice; it can still be named by hand.
		if p.Ext != 0 {
			continue
		}
		found = append(found, p)
	}
	return pick(found, func(p Platform) Version {
		return Version{Parts: []int{p.API, p.Minor}}
	})
}

// parsePlatform reads "android-30", "android-36.1" and "android-34-ext11".
func parsePlatform(name string) (Platform, bool) {
	rest, ok := strings.CutPrefix(name, "android-")
	if !ok {
		return Platform{}, false
	}
	p := Platform{Name: name}
	if base, ext, found := strings.Cut(rest, "-ext"); found {
		n, err := strconv.Atoi(ext)
		if err != nil {
			return Platform{}, false
		}
		p.Ext, rest = n, base
	}
	major, minor, _ := strings.Cut(rest, ".")
	n, err := strconv.Atoi(major)
	if err != nil {
		// A preview platform is named by a letter — "android-VanillaIceCream"
		// — and has no API level to compare, so it is skipped.
		return Platform{}, false
	}
	p.API = n
	if minor != "" {
		p.Minor, _ = strconv.Atoi(minor)
	}
	return p, true
}

func findNDK(root string) (NDK, bool) {
	var found []NDK
	add := func(dir string) {
		if dir == "" {
			return
		}
		n, ok := readNDK(dir)
		if ok {
			found = append(found, n)
		}
	}
	for _, env := range []string{"ANDROID_NDK_HOME", "ANDROID_NDK_ROOT", "NDK_ROOT"} {
		add(os.Getenv(env))
	}
	if len(found) > 0 {
		return found[0], true // an environment that names one means it
	}
	for _, name := range subdirs(filepath.Join(root, "ndk")) {
		add(filepath.Join(root, "ndk", name))
	}
	add(filepath.Join(root, "ndk-bundle"))
	return pick(found, func(n NDK) Version { return n.Version })
}

// readNDK reads an NDK's version out of source.properties and finds the
// prebuilt toolchain built for this machine.
func readNDK(dir string) (NDK, bool) {
	b, err := os.ReadFile(filepath.Join(dir, "source.properties"))
	if err != nil {
		return NDK{}, false
	}
	n := NDK{Dir: dir}
	for line := range strings.Lines(string(b)) {
		if rev, ok := strings.CutPrefix(strings.TrimSpace(line), "Pkg.Revision"); ok {
			if _, v, found := strings.Cut(rev, "="); found {
				n.Version = ParseVersion(strings.TrimSpace(v))
			}
		}
	}
	prebuilt := filepath.Join(dir, "toolchains", "llvm", "prebuilt")
	hosts := subdirs(prebuilt)
	if len(hosts) == 0 {
		return NDK{}, false
	}
	// Take the directory whose name starts with this OS. The NDK has shipped
	// darwin-x86_64 on Apple Silicon and darwin-arm64 in different revisions,
	// so the exact name is read rather than assumed, and any host directory
	// is better than refusing to build.
	n.Host = hosts[0]
	for _, h := range hosts {
		if strings.HasPrefix(h, runtime.GOOS) {
			n.Host = h
			break
		}
	}
	return n, n.Version.Major() > 0
}

// findBundletool looks where someone would have put it. It is never an
// error not to find it: it is not part of the SDK and nothing needs it to
// build.
func findBundletool(root string) string {
	if p := os.Getenv("BUNDLETOOL"); p != "" && exists(p) {
		return p
	}
	home, _ := os.UserHomeDir()
	for _, p := range []string{
		filepath.Join(root, "bundletool.jar"),
		filepath.Join(home, "bundletool.jar"),
		filepath.Join(home, ".local", "share", "bundletool.jar"),
		"/usr/local/lib/bundletool.jar",
	} {
		if exists(p) {
			return p
		}
	}
	if p, err := exec.LookPath("bundletool"); err == nil {
		return p
	}
	return ""
}

func findJDK() (JDK, bool) {
	if home := os.Getenv("JAVA_HOME"); home != "" {
		javac := filepath.Join(home, "bin", exe("javac"))
		if _, err := os.Stat(javac); err == nil {
			return JDK{Dir: home, Javac: javac}, true
		}
	}
	javac, err := exec.LookPath("javac")
	if err != nil {
		return JDK{}, false
	}
	// Follow the symlink an alternatives system leaves in /usr/bin, so the
	// home we report is the real one and not /usr.
	if real, err := filepath.EvalSymlinks(javac); err == nil {
		javac = real
	}
	return JDK{Dir: filepath.Dir(filepath.Dir(javac)), Javac: javac}, true
}

// subdirs lists a directory's subdirectories, sorted, and gives back nothing
// at all when the directory does not exist.
func subdirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		// A symlinked SDK directory is common enough to be worth following.
		if e.IsDir() {
			names = append(names, e.Name())
			continue
		}
		if e.Type()&os.ModeSymlink != 0 {
			if st, err := os.Stat(filepath.Join(dir, e.Name())); err == nil && st.IsDir() {
				names = append(names, e.Name())
			}
		}
	}
	return names
}

// exe adds the extension a program has on this machine.
func exe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

// script adds the extension a *wrapper* has: the SDK ships d8 and apksigner
// as shell scripts everywhere and as .bat files on Windows, which is not the
// same rule as a compiled tool like aapt2.
func script(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".bat"
	}
	return name
}
