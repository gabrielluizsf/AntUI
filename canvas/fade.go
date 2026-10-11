package canvas

// Fade replaces the colour's alpha, keeping its RGB.
func Fade(c Color, alpha int) Color {
	return (c & 0x00FFFFFF) | Color(uint32(clampInt(alpha, 0, 255))<<24)
}
