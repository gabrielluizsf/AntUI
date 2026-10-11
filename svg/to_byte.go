package svg

// toByte rounds a channel to the byte it becomes, holding it inside 0-255 so a
// value outside that range is clamped rather than wrapped around.
func toByte(v float64) uint8 { return byte(clampFloat(v, 0, 255) + 0.5) }
