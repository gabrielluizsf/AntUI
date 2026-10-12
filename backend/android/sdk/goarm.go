package sdk

// GOARM is the ARM version to build for, and is only ever set for the 32-bit
// ARM ABI — armeabi-v7a means ARMv7 with hardware float, which is GOARM=7.
func (a ABI) GOARM() string {
	if a == Arm {
		return "7"
	}
	return ""
}
