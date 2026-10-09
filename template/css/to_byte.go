package css

// clampByte rounds a 0-1 component to a byte.
func toByte(v float64) uint8 {
	return uint8(clamp01(v)*255 + 0.5)
}
