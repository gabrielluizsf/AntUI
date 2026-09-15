package sdk

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

// Jar is the android.jar the Java shim compiles against.
func (p Platform) Jar() string { return filepath.Join(p.Dir, "android.jar") }

// Target is the API level to write into a manifest as targetSdkVersion. It
// is the major number: the minor of an android-36.1 is a mid-year addition
// to the same level and is not what targetSdkVersion holds.
func (p Platform) Target() int { return p.API }

// String is the platform as a person would say it: "android-36".
func (p Platform) String() string {
	if p.Name != "" {
		return p.Name
	}
	return "android-" + strconv.Itoa(p.API)
}

// The four tools an APK is built with. Two of them, d8 and apksigner, are
// wrapper scripts rather than programs, which only shows on Windows.
func (b BuildTools) Aapt2() string { return filepath.Join(b.Dir, exe("aapt2")) }

// Zipalign is the tool that aligns a package so the loader can map what
// is inside it.
func (b BuildTools) Zipalign() string { return filepath.Join(b.Dir, exe("zipalign")) }

// D8 is the tool that turns Java class files into a dex.
func (b BuildTools) D8() string { return filepath.Join(b.Dir, script("d8")) }

// Apksigner is the tool that signs an APK, and the only one that can
// write the v2 and v3 signatures a modern Android insists on.
func (b BuildTools) Apksigner() string { return filepath.Join(b.Dir, script("apksigner")) }

// Aidl is the tool that compiles an interface definition. Nothing here
// uses it; it is listed because a build-tools install without it is a
// broken one.
func (b BuildTools) Aidl() string { return filepath.Join(b.Dir, exe("aidl")) }

// String is the build-tools version as a person would say it.
func (b BuildTools) String() string { return "build-tools " + b.Version.String() }

// Adb, and the rest of what lives outside build-tools.
func (t *Toolchain) Adb() string { return filepath.Join(t.Root, "platform-tools", exe("adb")) }

// Emulator is the emulator, which is a separate download from the SDK.
func (t *Toolchain) Emulator() string { return filepath.Join(t.Root, "emulator", exe("emulator")) }

// Sdkmanager is what installs and updates everything else.
func (t *Toolchain) Sdkmanager() string { return t.cmdlineTool("sdkmanager") }

// Avdmanager is what makes and lists the virtual devices.
func (t *Toolchain) Avdmanager() string { return t.cmdlineTool("avdmanager") }

func (t *Toolchain) cmdlineTool(name string) string {
	if runtime.GOOS == "windows" {
		name += ".bat"
	}
	return filepath.Join(t.Root, "cmdline-tools", "latest", "bin", name)
}

// Bin is the NDK's toolchain directory, where clang and the llvm tools are.
func (n NDK) Bin() string {
	return filepath.Join(n.Dir, "toolchains", "llvm", "prebuilt", n.Host, "bin")
}

// Clang is the compiler that builds for one ABI against one API level. The
// NDK ships a wrapper per (target, level) pair rather than one compiler
// taking a flag, so the API level is baked into the name and a level with no
// wrapper is an error here rather than a confusing failure later.
func (n NDK) Clang(abi ABI, api int) (string, error) {
	if !abi.Valid() {
		return "", fmt.Errorf("sdk: %q is not an Android ABI", abi)
	}
	name := abi.Triple() + strconv.Itoa(api) + "-clang"
	if runtime.GOOS == "windows" {
		name += ".cmd"
	}
	path := filepath.Join(n.Bin(), name)
	if !exists(path) {
		return "", fmt.Errorf("sdk: NDK %s has no compiler for %s at API %d (%s)",
			n.Version, abi, api, name)
	}
	return path, nil
}

// Tool is one of the llvm programs beside clang — llvm-strip, llvm-objcopy,
// llvm-readelf — named without its "llvm-" prefix.
func (n NDK) Tool(name string) string { return filepath.Join(n.Bin(), exe("llvm-"+name)) }

// String is the NDK version as a person would say it.
func (n NDK) String() string { return "NDK " + n.Version.String() }

// Beta reports whether this NDK is a preview. Google publishes betas into the
// same channel as releases, so an install can quietly be one; a build still
// works, but it is not what to ship a store bundle from.
func (n NDK) Beta() bool { return !n.Version.Release() }
