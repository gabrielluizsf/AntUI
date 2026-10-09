package css

// clampByte rounds a 0-1 component to a byte.
func toByte(v float64) uint8 {
	return uint8(clamp(v)*255 + 0.5)
}
