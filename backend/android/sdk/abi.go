package sdk

import "fmt"

// ABI is one of the machine architectures Android runs on, named the way the
// platform names it — these strings are the directory names inside an APK's
// lib/ folder, so they are spelt Android's way and not Go's.
type ABI string

// The four ABIs that still exist. ARMv5 and MIPS are gone.
const (
	Arm64 ABI = "arm64-v8a"   // every phone sold since 2016
	Arm   ABI = "armeabi-v7a" // 32-bit ARM, still on cheap and old devices
	X64   ABI = "x86_64"      // the emulator, and Chromebooks
	X86   ABI = "x86"         // 32-bit x86, effectively only old emulators
)

// Ship is what goes to the store: 64-bit ARM is required, and 32-bit ARM is
// what keeps the app installable on a phone from before 2016. The x86 pair
// is for the emulator and is deliberately not here — including it in a
// bundle costs every user the download of a binary no phone can run.
var Ship = []ABI{Arm64, Arm}

// All is every ABI, newest first.
var All = []ABI{Arm64, Arm, X64, X86}

// ParseABI reads an ABI name, also accepting the Go architecture that maps
// onto it so a caller can pass GOARCH through unchanged.
func ParseABI(s string) (ABI, error) {
	switch s {
	case "arm64-v8a", "arm64":
		return Arm64, nil
	case "armeabi-v7a", "armeabi", "arm":
		return Arm, nil
	case "x86_64", "amd64":
		return X64, nil
	case "x86", "386":
		return X86, nil
	}
	return "", fmt.Errorf("sdk: %q is not an Android ABI", s)
}

// GOARCH is the Go architecture that builds for this ABI.
func (a ABI) GOARCH() string {
	switch a {
	case Arm64:
		return "arm64"
	case Arm:
		return "arm"
	case X64:
		return "amd64"
	case X86:
		return "386"
	}
	return ""
}

// GOARM is the ARM version to build for, and is only ever set for the 32-bit
// ARM ABI — armeabi-v7a means ARMv7 with hardware float, which is GOARM=7.
func (a ABI) GOARM() string {
	if a == Arm {
		return "7"
	}
	return ""
}

// Triple is the target the NDK's clang is named after. Note that 32-bit ARM
// is "armv7a-linux-androideabi" here and "arm-linux-androideabi" everywhere
// binutils is involved; this is the compiler's spelling, which is the one the
// clang wrapper scripts use.
func (a ABI) Triple() string {
	switch a {
	case Arm64:
		return "aarch64-linux-android"
	case Arm:
		return "armv7a-linux-androideabi"
	case X64:
		return "x86_64-linux-android"
	case X86:
		return "i686-linux-android"
	}
	return ""
}

// Bits is 64 or 32. The Play Store has refused 32-bit-only uploads since
// August 2019, so a bundle has to carry at least one 64-bit ABI.
func (a ABI) Bits() int {
	switch a {
	case Arm64, X64:
		return 64
	}
	return 32
}

// Valid reports whether the ABI is one this package knows.
func (a ABI) Valid() bool { return a.GOARCH() != "" }

// String is the ABI as Android spells it, which is also the directory
// name inside a package.
func (a ABI) String() string { return string(a) }
