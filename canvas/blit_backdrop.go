package canvas

import "math"

// BlitBackdrop lays src over this canvas sampling through the forward map:
// destination point p reads src at m(p), where BlitMatrix would read at
// m⁻¹(p). The destination region is fixed here — every point inside it is
// pulled from src — which is the direction a transformed element needs when
// it plants the window content sitting under its box into its recording
// layer before a backdrop-filter looks at it. The same nearest-versus-
// bilinear split as BlitMatrix applies.
func (cv *Canvas) BlitBackdrop(src *Canvas, m Matrix, region Area) {
	if src == nil || src.Width <= 0 || src.Height <= 0 {
		return
	}
	x0 := max(region.X, cv.Clip.X)
	y0 := max(region.Y, cv.Clip.Y)
	x1 := min(region.X+region.Width, cv.Clip.X+cv.Clip.Width)
	y1 := min(region.Y+region.Height, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	axis := m.B == 0 && m.C == 0
	for dy := y0; dy < y1; dy++ {
		for dx := x0; dx < x1; dx++ {
			sx, sy := m.Map(float64(dx)+0.5, float64(dy)+0.5)
			var col Color
			if axis {
				ix, iy := int(math.Floor(sx)), int(math.Floor(sy))
				if ix < 0 || iy < 0 || ix >= src.Width || iy >= src.Height {
					continue
				}
				col = src.At(ix, iy)
			} else {
				col = bilinearAt(src, 0, 0, src.Width, src.Height, sx, sy)
			}
			if col == 0 {
				continue
			}
			cv.Pixel(dx, dy, col)
		}
	}
}
