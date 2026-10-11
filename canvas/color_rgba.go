package canvas

// RGBA implements image/color.Color, so a canvas.Color can be handed to
// anything in the standard library that draws. The channels are
// alpha-premultiplied and scaled to 16 bits, as that interface requires.
func (c Color) RGBA() (r, g, b, a uint32) {
	al := uint32(c.A())
	al |= al << 8
	r = uint32(c.R()) * al / 0xFF
	g = uint32(c.G()) * al / 0xFF
	b = uint32(c.B()) * al / 0xFF
	return r, g, b, al
}
