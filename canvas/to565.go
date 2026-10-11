package canvas

// To565 packs a colour into 16 bits. It rounds rather than truncating.
func To565(c Color) uint16 {
	r5 := (uint32(c.R())*31 + 127) / 255
	g6 := (uint32(c.G())*63 + 127) / 255
	b5 := (uint32(c.B())*31 + 127) / 255
	return uint16(r5<<11 | g6<<5 | b5)
}
