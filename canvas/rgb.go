package canvas

// RGB builds an opaque colour.
func RGB(r, g, b uint8) Color {
	return Color(0xFF000000 | uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}
