package font

import "math"

// line adds one edge, in pixels, with y downwards.
//
// This is the whole rasteriser. Everything else is arranging to call it.
func (r *raster) line(x0, y0, x1, y1 float64) {
	if r.w == 0 || y0 == y1 {
		return
	}
	// Going up counts one way and going down the other, which is what makes
	// a hole inside a letter a hole: the outer loop and the inner loop wind
	// in opposite directions and cancel.
	dir := 1.0
	if y0 > y1 {
		dir = -1
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	dxdy := (x1 - x0) / (y1 - y0)

	yStart := max(int(math.Floor(y0)), 0)
	yEnd := min(int(math.Ceil(y1)), r.h)

	for y := yStart; y < yEnd; y++ {
		// The part of this scanline the edge actually spans.
		top := max(y0, float64(y))
		bottom := min(y1, float64(y+1))
		height := bottom - top
		if height <= 0 {
			continue
		}
		xTop := x0 + dxdy*(top-y0)
		xBottom := x0 + dxdy*(bottom-y0)
		r.span(y, xTop, xBottom, height*dir)
	}
}
