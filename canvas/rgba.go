package canvas

// RGBAParams groups components of an RGBA colour.
type RGBAParams struct {
	R, G, B, A uint8
}

// RGBA packs the colour components.
func RGBA(p RGBAParams) Color {
	return Color(uint32(p.A)<<24 | uint32(p.R)<<16 | uint32(p.G)<<8 | uint32(p.B))
}
