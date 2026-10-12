package sdk

// The four ABIs that still exist. ARMv5 and MIPS are gone.
const (
	Arm64 ABI = "arm64-v8a"   // every phone sold since 2016
	Arm   ABI = "armeabi-v7a" // 32-bit ARM, still on cheap and old devices
	X64   ABI = "x86_64"      // the emulator, and Chromebooks
	X86   ABI = "x86"         // 32-bit x86, effectively only old emulators
)
