package canvas

// Blend composites src over dst, honouring the alpha of src. The result is
// always opaque: this is what lands in a framebuffer, and a framebuffer has
// nothing behind it to show through.
func Blend(dst, src Color) Color {
	a := uint32(src>>24) & 0xFF
	if a == 255 {
		return src | 0xFF000000
	}
	if a == 0 {
		return dst
	}
	inv := 255 - a
	r := ((uint32(src>>16)&0xFF)*a + (uint32(dst>>16)&0xFF)*inv + 127) / 255
	g := ((uint32(src>>8)&0xFF)*a + (uint32(dst>>8)&0xFF)*inv + 127) / 255
	b := ((uint32(src)&0xFF)*a + (uint32(dst)&0xFF)*inv + 127) / 255
	return Color(0xFF000000 | r<<16 | g<<8 | b)
}
