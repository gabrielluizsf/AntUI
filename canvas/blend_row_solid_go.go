package canvas

func blendRowSolidGo(dst []Color, c Color) {
	for i := range dst {
		dst[i] = Blend(dst[i], c)
	}
}
