package css

// hwbToRGB turns hue/whiteness/blackness percentages into sRGB. When white
// and black fill the circle a neutral grey emerges; otherwise the pure hue is
// scaled and mixed toward white and black.
func hwbToRGB(h, w, b float64) [3]uint8 {
	w /= 100
	b /= 100
	if w+b >= 1 {
		g := uint8(w/(w+b)*255 + 0.5)
		return [3]uint8{g, g, g}
	}
	// The hue at maximum chroma, then constricted by the white and black.
	hue := hslToRGB(h, 100, 50)
	f := 1 - w - b
	return [3]uint8{
		uint8(float64(hue[0])*f + w*255 + 0.5),
		uint8(float64(hue[1])*f + w*255 + 0.5),
		uint8(float64(hue[2])*f + w*255 + 0.5),
	}
}
