package sdk

// Bits is 64 or 32. The Play Store has refused 32-bit-only uploads since
// August 2019, so a bundle has to carry at least one 64-bit ABI.
func (a ABI) Bits() int {
	switch a {
	case Arm64, X64:
		return 64
	}
	return 32
}
