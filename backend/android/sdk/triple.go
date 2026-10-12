package sdk

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
