package canvas

// BlendOver composites src over dst honouring both alphas, so the result
// keeps the transparency the mix leaves behind. Blend treats dst as an opaque
// framebuffer and always answers opaque — right for painting the window,
// wrong for stamping a soft shadow on a recording layer that is still
// transparent. Over an opaque dst the two agree pixel for pixel.
func BlendOver(dst, src Color) Color {
	a := int((src >> 24) & 0xFF)
	if a == 255 {
		return src
	}
	da := int((dst >> 24) & 0xFF)
	if a == 0 {
		return dst
	}
	dstW := da * (255 - a) / 255
	out := a + dstW
	if out == 0 {
		return 0
	}
	ch := func(s, d int) uint8 {
		return uint8((s*a + d*dstW + out/2) / out)
	}
	return Color(uint32(out)<<24 |
		uint32(ch(int((src>>16)&0xFF), int((dst>>16)&0xFF)))<<16 |
		uint32(ch(int((src>>8)&0xFF), int((dst>>8)&0xFF)))<<8 |
		uint32(ch(int(src&0xFF), int(dst&0xFF))))
}
