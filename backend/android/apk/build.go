package apk

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// Result is what a build produced.
type Result struct {
	// APK is the file written.
	APK string
	// Libs is the compiled library for each ABI, in the temporary directory
	// the build used. They are gone by the time Build returns unless Keep
	// was set.
	Libs map[sdk.ABI]string
	// Size is the finished package, in bytes.
	Size int64
	// Symbols is the archive of symbol files written beside the package, or
	// empty when the build was told not to strip and the symbols are
	// therefore still in the shipped library. See [WriteSymbols].
	Symbols string
}

// Build compiles the Go package and produces a signed, installable APK.
//
// It is the four steps below run in order, in a temporary directory that is
// removed afterwards. Nothing is written outside it except the APK itself.
func Build(tc *sdk.Toolchain, cfg Config) (*Result, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	if tc.NDK.Dir == "" {
		return nil, fmt.Errorf("apk: no NDK; a Go app is native code and cannot be built without one")
	}
	if tc.BuildTools.Dir == "" || tc.Platform.Dir == "" {
		return nil, fmt.Errorf("apk: the SDK has no build-tools or no platform installed")
	}

	work, err := os.MkdirTemp("", "antuiapk-")
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	defer os.RemoveAll(work)

	libs, err := BuildLibs(tc, &cfg, work)
	if err != nil {
		return nil, err
	}
	base := filepath.Join(work, "base.apk")
	if err := Link(tc, &cfg, base); err != nil {
		return nil, err
	}
	packed := filepath.Join(work, "packed.apk")
	if err := Assemble(base, libs, cfg.Lib, packed, &cfg); err != nil {
		return nil, err
	}
	if err := Sign(tc, &cfg, packed, cfg.Out); err != nil {
		return nil, err
	}

	st, err := os.Stat(cfg.Out)
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	res := &Result{APK: cfg.Out, Libs: libs, Size: st.Size()}
	// The symbols go out beside the package, because the moment they are
	// made is the only moment they can be kept.
	syms := symbolsPath(&cfg)
	if wrote, err := WriteSymbols(libs, cfg.Lib, syms); err != nil {
		return nil, err
	} else if wrote {
		res.Symbols = syms
	}
	return res, nil
}

// BuildLibs compiles the Go package once per ABI, writing each shared
// library into dir. The compiler is the NDK's clang for that ABI at the
// app's minimum SDK — not the newest available, because a library built
// against a newer platform's headers can reference symbols the oldest device
// it claims to support does not have, and that failure happens at load time
// on a user's phone rather than here.
func BuildLibs(tc *sdk.Toolchain, cfg *Config, dir string) (map[sdk.ABI]string, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	goBin, err := goCommand()
	if err != nil {
		return nil, err
	}
	overlay, err := overlayFor(cfg, dir)
	if err != nil {
		return nil, err
	}

	libs := make(map[sdk.ABI]string, len(cfg.ABIs))
	for _, abi := range cfg.ABIs {
		cc, err := tc.NDK.Clang(abi, cfg.MinSDK)
		if err != nil {
			return nil, fmt.Errorf("apk: %w", err)
		}
		out := filepath.Join(dir, string(abi), "lib"+cfg.Lib+".so")
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, fmt.Errorf("apk: %w", err)
		}

		env := append(os.Environ(),
			"GOOS=android",
			"GOARCH="+abi.GOARCH(),
			"CGO_ENABLED=1",
			"CC="+cc,
			"CXX="+strings.TrimSuffix(cc, "clang")+"clang++",
		)
		if arm := abi.GOARM(); arm != "" {
			env = append(env, "GOARM="+arm)
		}

		// Nothing is stripped at the link. The library is built whole and
		// then split into what ships and what reads a crash — see
		// [keepSymbols], which explains why that is not the same as -s -w
		// even though the shipped library comes out identical in size.
		args := []string{"build", "-buildmode=c-shared", "-trimpath",
			"-overlay", overlay, "-o", out}
		if len(cfg.Tags) > 0 {
			args = append(args, "-tags", strings.Join(cfg.Tags, ","))
		}
		if cfg.Ldflags != "" {
			args = append(args, "-ldflags", cfg.Ldflags)
		}
		args = append(args, ".")

		if _, err := runIn(cfg.Dir, env, goBin, args...); err != nil {
			return nil, fmt.Errorf("apk: building for %s: %w", abi, err)
		}
		if *cfg.Strip {
			if err := keepSymbols(tc, out); err != nil {
				return nil, err
			}
		}
		libs[abi] = out
	}
	return libs, nil
}

// Link turns the manifest into a base APK with aapt2. The result holds the
// manifest in binary form and a resource table, and nothing else — no code
// and no libraries.
func Link(tc *sdk.Toolchain, cfg *Config, out string) error {
	return link(tc, cfg, out, false)
}

// link is Link and LinkProto, which differ by one flag and nothing else.
func link(tc *sdk.Toolchain, cfg *Config, out string, proto bool) error {
	if err := cfg.check(); err != nil {
		return err
	}
	dir := filepath.Dir(out)
	manifest := filepath.Join(dir, "AndroidManifest.xml")
	if err := os.WriteFile(manifest, []byte(cfg.Manifest()), 0o644); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	// Resources, when there are any. aapt2 works in two steps and they are
	// not interchangeable: compile turns a directory of files into a flat
	// archive of compiled ones, and link takes those archives and the
	// manifest and makes the package. Handing link a directory of PNGs does
	// nothing at all.
	var compiled string
	if cfg.Icon != "" {
		var err error
		// The icons and their compiled form depend on nothing but the
		// picture, so this is usually a cache read — see [compiledIcons].
		compiled, err = compiledIcons(tc, dir, cfg.Icon, cfg.IconBackground)
		if err != nil {
			return err
		}
	}

	args := []string{"link",
		"-I", tc.Platform.Jar(),
		"--manifest", manifest,
		"-o", out,
		"--min-sdk-version", fmt.Sprint(cfg.MinSDK),
		"--target-sdk-version", fmt.Sprint(cfg.TargetSDK),
		"--version-code", fmt.Sprint(cfg.VersionCode),
		"--version-name", cfg.VersionName,
	}
	if cfg.Assets != "" {
		args = append(args, "-A", cfg.Assets)
		// Anything already compressed goes in stored, so that it can be
		// mapped rather than inflated on every read.
		for _, ext := range cfg.Store {
			args = append(args, "-0", ext)
		}
	}
	// The compiled resources go on the end as a positional argument, not
	// behind a flag. -R exists and means something else: it is for *overlay*
	// resources, and passing the ordinary ones to it fails with "failed
	// parsing overlays", which does not suggest what to do about it.
	if proto {
		// The manifest and the resource table come out as protocol buffers
		// rather than binary XML and an arsc, which is what a bundle holds.
		args = append(args, "--proto-format")
	}
	if compiled != "" {
		args = append(args, compiled)
	}
	if _, err := run(tc.BuildTools.Aapt2(), args...); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}

// Assemble copies the base APK and adds the libraries — and the Java shim,
// when the config asks for one — to it.
//
// The libraries go in **stored, not deflated**. That is not an optimisation:
// the manifest says extractNativeLibs="false", so the loader maps the .so
// straight out of the APK, and a compressed entry cannot be mapped. The
// alignment that mapping also needs is [Sign]'s job.
//
// Everything already in the base APK is copied byte for byte, keeping
// whatever compression aapt2 chose — except resources.arsc, which since API
// 30 must also be stored, and is decompressed here if it is not.
func Assemble(base string, libs map[sdk.ABI]string, libName, out string, cfg *Config) error {
	zr, err := zip.OpenReader(base)
	if err != nil {
		return fmt.Errorf("apk: reading %s: %w", base, err)
	}
	defer zr.Close()

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)

	for _, e := range zr.File {
		if e.Name == "resources.arsc" && e.Method != zip.Store {
			if err := storeEntry(zw, e); err != nil {
				return err
			}
			continue
		}
		if err := copyRaw(zw, e); err != nil {
			return err
		}
	}
	for abi, lib := range libs {
		name := path.Join("lib", string(abi), "lib"+libName+".so")
		if err := addStored(zw, name, lib); err != nil {
			return err
		}
	}
	// The shim, when the manifest says there is one. It is deflated rather
	// than stored: a dex is read by the runtime rather than mapped, so
	// nothing needs it uncompressed, and it compresses well.
	if cfg != nil && cfg.Shim != nil && *cfg.Shim {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: "classes.dex", Method: zip.Deflate})
		if err != nil {
			return fmt.Errorf("apk: classes.dex: %w", err)
		}
		if _, err := w.Write(shim.Dex); err != nil {
			return fmt.Errorf("apk: classes.dex: %w", err)
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return f.Close()
}

// copyRaw moves an entry across without recompressing it.
func copyRaw(zw *zip.Writer, e *zip.File) error {
	rc, err := e.OpenRaw()
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	h := e.FileHeader
	w, err := zw.CreateRaw(&h)
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	return nil
}

// storeEntry rewrites a compressed entry as a stored one.
func storeEntry(zw *zip.Writer, e *zip.File) error {
	rc, err := e.Open()
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	defer rc.Close()
	w, err := zw.CreateHeader(&zip.FileHeader{Name: e.Name, Method: zip.Store})
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	return nil
}

func addStored(zw *zip.Writer, name, from string) error {
	src, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer src.Close()
	w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
	if err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	return nil
}

// Sign aligns the package and signs it.
//
// The order is not free: zipalign moves entries, so signing first would
// invalidate the signature. apksigner knows this and refuses to sign
// something it would then have to be run over again.
//
// The alignment asked for is 16 KB, not the 4 KB that was standard for
// years. Android 15 introduced devices with 16 KB memory pages, and the
// Play Store requires an app to work on them; a library aligned to 4 KB
// cannot be mapped on such a device at all.
func Sign(tc *sdk.Toolchain, cfg *Config, in, out string) error {
	if err := cfg.check(); err != nil {
		return err
	}
	return SignAPK(tc, cfg.Keystore, cfg.MinSDK, in, out)
}

// SignAPK aligns and signs a package that is already built, with no config
// around it. It is what [Sign] does, and what "antuiapk sign" is: an APK
// that came from somewhere else — an unsigned one, or one to be re-signed
// with a release key — has no Go package behind it to describe.
func SignAPK(tc *sdk.Toolchain, k Keystore, minSDK int, in, out string) error {
	if minSDK == 0 {
		minSDK = sdk.MinSDK
	}
	// The zero value is the debug key, which is what [Config.check] would
	// have done for a build. Signing on its own goes through neither.
	k.fill()
	if err := k.Ensure(tc); err != nil {
		return err
	}
	aligned := in + ".aligned"
	if _, err := run(tc.BuildTools.Zipalign(), "-P", "16", "-f", "4", in, aligned); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(mustAbs(out)), 0o755); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	_, err := run(tc.BuildTools.Apksigner(), "sign",
		"--ks", k.Path,
		"--ks-key-alias", k.Alias,
		"--ks-pass", "pass:"+k.StorePass,
		"--key-pass", "pass:"+k.KeyPass,
		"--min-sdk-version", fmt.Sprint(minSDK),
		// v4 is off. It is a detached signature in a second file beside the
		// APK, used only to make "adb install" incremental, and leaving a
		// stray .idsig next to whatever the caller asked for is a worse
		// surprise than a slower install.
		"--v4-signing-enabled", "false",
		"--out", out,
		aligned,
	)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}

func mustAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

// goCommand finds the go tool: on PATH first, and in GOROOT when the build
// is running somewhere PATH was never set up, which is most CI.
func goCommand() (string, error) {
	if p, err := exec.LookPath("go"); err == nil {
		return p, nil
	}
	p := filepath.Join(runtime.GOROOT(), "bin", "go")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("apk: no go tool on PATH")
}
