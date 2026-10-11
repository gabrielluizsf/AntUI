package canvas

import "math"

// BlitMatrix composites src over this canvas through the affine matrix m,
// which maps src pixel coordinates to this canvas's. region bounds the part of
// src worth walking — only its transformed footprint is drawn — so a
// full-window recording layer costs only the pixels its widget covered.
//
// Where m only translates, mirrors or scales the axis-aligned path samples
// the nearest pixel, keeping such transforms crisp. A rotation or a shear
// switches to bilinear sampling: the four neighbours of the mapped point are
// blended in premultiplied space, with a largely transparent neighbour
// darkening nothing.
func (cv *Canvas) BlitMatrix(src *Canvas, m Matrix, region Area) {
	if src == nil || src.Width <= 0 || src.Height <= 0 {
		return
	}
	inv, ok := m.Inverse()
	if !ok {
		return
	}

	// The part of src that is worth mapping, clamped to the canvas.
	x0 := max(region.X, 0)
	y0 := max(region.Y, 0)
	x1 := min(region.X+region.Width, src.Width)
	y1 := min(region.Y+region.Height, src.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}

	d0x, d1x, d0y, d1y, ok := blitMatrixBox(cv, m, x0, y0, x1, y1)
	if !ok {
		return
	}

	axis := m.B == 0 && m.C == 0
	for dy := d0y; dy < d1y; dy++ {
		for dx := d0x; dx < d1x; dx++ {
			sx, sy := inv.Map(float64(dx)+0.5, float64(dy)+0.5)
			var col Color
			if axis {
				ix, iy := int(math.Floor(sx)), int(math.Floor(sy))
				if ix < x0 || iy < y0 || ix >= x1 || iy >= y1 {
					continue
				}
				col = src.At(ix, iy)
			} else {
				col = bilinearAt(src, x0, y0, x1, y1, sx, sy)
			}
			if col == 0 {
				continue
			}
			cv.Pixel(dx, dy, col)
		}
	}
}
