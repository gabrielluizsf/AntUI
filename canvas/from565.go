package canvas

// From565 unpacks a 16-bit colour, repeating the top bits into the bottom.
func From565(packed uint16) Color {
	r5 := uint32(packed>>11) & 0x1F
	g6 := uint32(packed>>5) & 0x3F
	b5 := uint32(packed) & 0x1F
	r := r5<<3 | r5>>2
	g := g6<<2 | g6>>4
	b := b5<<3 | b5>>2
	return Color(0xFF000000 | r<<16 | g<<8 | b)
}
