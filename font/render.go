package font

import "math"

// Render draws a glyph's contours into a mask.
//
// scale turns font units into pixels, and the origin is the glyph's own: x to
// the right of the pen, y **up** from the baseline, which is the opposite of
// the canvas and is why y is flipped on the way in.
func Render(contours []Contour, scale float64) Mask {
	if len(contours) == 0 {
		return Mask{}
	}
	// The box the glyph actually occupies, so the mask is only as big as it
	// has to be and a page of text is not a page of empty pixels.
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, c := range contours {
		for i := range c.X {
			minX, maxX = math.Min(minX, c.X[i]), math.Max(maxX, c.X[i])
			minY, maxY = math.Min(minY, c.Y[i]), math.Max(maxY, c.Y[i])
		}
	}
	left := int(math.Floor(minX * scale))
	right := int(math.Ceil(maxX * scale))
	bottom := int(math.Floor(minY * scale))
	top := int(math.Ceil(maxY * scale))
	w, h := right-left+1, top-bottom+1
	if w <= 0 || h <= 0 || w > 4096 || h > 4096 {
		return Mask{}
	}

	r := newRaster(w, h)
	for _, c := range contours {
		Emit(r, c, scale, float64(left), float64(top))
	}
	return Mask{W: w, H: h, Left: left, Top: top, Alpha: r.mask(), Stride: w}
}
