// Package shim holds the only Java in this library, and the compiled form of
// it that ships inside an APK.
//
// Three things Android delivers by overriding a method on the activity —
// the result of a permission request, the result of another activity, and an
// intent arriving at one already running — cannot be reached from native
// code, because there is no way to override a Java method from C. So there
// is one small subclass of NativeActivity that overrides those three and
// hands each one's arguments straight to Go.
//
// The dex is **committed**, not built. An app made with this library needs
// no JDK and no Java toolchain; only changing the shim itself does, and
// [Build] is what does it. d8 is deterministic, so the committed file can be
// checked against the source rather than trusted.
package shim

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Dex is the compiled shim, ready to be put in an APK as classes.dex.
//
//go:embed classes.dex
var Dex []byte

// Sources is the Java it was compiled from — every file, so that [Build] can
// reproduce the dex without knowing where the repository is.
//
//go:embed *.java
var Sources embed.FS

// ProGuard is what R8 must be told to leave alone.
//
// AntUI's own pipeline never runs R8: antuiapk compiles the shim with javac
// and dexes it with d8, and neither shrinks anything. This is for an app
// that merges the shim into a Gradle build, where R8 is on for every release
// build and would otherwise remove most of the shim as dead code — nothing
// in it is called from Java, so to a shrinker none of it is reachable.
//
// The failure it prevents is a quiet one: the build succeeds, the app
// installs, and a JNI call throws NoSuchMethodError the first time that code
// path runs.
//
//	os.WriteFile("proguard-rules.pro", shim.ProGuard, 0o644)
//
//go:embed proguard-rules.pro
var ProGuard []byte

// Package is the Java package every shim class is in. Native code finds them
// by name, so the name is part of the interface rather than an arrangement
// of files.
const Package = "dev.antui"

// Activity is the class an app's manifest points at when the shim is in, and
// Receiver is the one a notification's action buttons broadcast to.
const (
	Activity = "dev.antui.AntuiActivity"
	Receiver = "dev.antui.AntuiReceiver"
	// Provider is the one that hands a file to another app, and
	// ProviderSuffix is what its authority is the app's package plus.
	Provider       = "dev.antui.AntuiFileProvider"
	ProviderSuffix = ".antuifiles"
)

// SourceLevel is the Java version the shim is compiled to. It is old on
// purpose: nothing here needs anything newer, and the older the class file
// the wider the range of build tools that will dex it.
const SourceLevel = "8"

// Build compiles the shim from [Source] and returns the dex.
//
// It needs a JDK and the SDK's build tools, which building an app does not —
// which is the whole reason the result is committed. Nothing is written
// outside the temporary directory it makes.
func Build(tc *sdk.Toolchain) ([]byte, error) {
	if tc.JDK.Javac == "" {
		return nil, fmt.Errorf("shim: rebuilding needs a JDK, and none was found")
	}
	if tc.Platform.Dir == "" {
		return nil, fmt.Errorf("shim: rebuilding needs an SDK platform to compile against")
	}
	if tc.BuildTools.Dir == "" {
		return nil, fmt.Errorf("shim: rebuilding needs the SDK build tools for d8")
	}

	work, err := os.MkdirTemp("", "antui-shim-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(work)

	names, err := fs.Glob(Sources, "*.java")
	if err != nil || len(names) == 0 {
		return nil, fmt.Errorf("shim: no Java to compile")
	}
	sort.Strings(names) // so the class files reach d8 in a fixed order
	srcs := make([]string, 0, len(names))
	for _, name := range names {
		b, err := Sources.ReadFile(name)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(work, name)
		if err := os.WriteFile(path, b, 0o644); err != nil {
			return nil, err
		}
		srcs = append(srcs, path)
	}
	classes := filepath.Join(work, "classes")
	if err := os.MkdirAll(classes, 0o755); err != nil {
		return nil, err
	}

	// Compiled against android.jar as the bootclasspath rather than the
	// JDK's own: the shim extends a platform class that no JDK has.
	args := []string{
		"-source", SourceLevel, "-target", SourceLevel,
		"-bootclasspath", tc.Platform.Jar(),
		"-nowarn",
		"-d", classes,
	}
	out, err := exec.Command(tc.JDK.Javac, append(args, srcs...)...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("shim: javac: %w\n%s", err, out)
	}

	var compiled []string
	err = filepath.WalkDir(classes, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(p, ".class") {
			compiled = append(compiled, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	if len(compiled) == 0 {
		return nil, fmt.Errorf("shim: javac produced no class files")
	}

	// d8 is deterministic given the same inputs in the same order, which is
	// what makes the committed dex checkable.
	sort.Strings(compiled)
	d8 := append([]string{"--lib", tc.Platform.Jar(), "--output", work}, compiled...)
	if out, err := exec.Command(tc.BuildTools.D8(), d8...).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("shim: d8: %w\n%s", err, out)
	}
	return os.ReadFile(filepath.Join(work, "classes.dex"))
}
