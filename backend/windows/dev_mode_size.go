//go:build windows

package windows

// DEVMODEW is large and its layout is fixed by the ABI rather than by any
// struct this package would want to write out. Only two of its fields matter
// here, so it is a byte buffer with the two offsets named.
const (
	devModeSize         = 220
	devModeSizeOffset   = 68
	devModeFrequencyOff = 184
)
