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

	// Destination bbox: the footprint of the region's corners.
	ax, ay := m.Map(float64(x0), float64(y0))
	bx, by := m.Map(float64(x1), float64(y0))
	cx, cy := m.Map(float64(x1), float64(y1))
	dx, dy := m.Map(float64(x0), float64(y1))
	minX := math.Min(math.Min(ax, bx), math.Min(cx, dx))
	maxX := math.Max(math.Max(ax, bx), math.Max(cx, dx))
	minY := math.Min(math.Min(ay, by), math.Min(cy, dy))
	maxY := math.Max(math.Max(ay, by), math.Max(cy, dy))

	d0x := max(int(math.Floor(minX)), cv.Clip.X)
	d1x := min(int(math.Ceil(maxX)), cv.Clip.X+cv.Clip.Width)
	d0y := max(int(math.Floor(minY)), cv.Clip.Y)
	d1y := min(int(math.Ceil(maxY)), cv.Clip.Y+cv.Clip.Height)
	if d0x >= d1x || d0y >= d1y {
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

// bilinearAt samples src bilinearly at the continuous point s,t, outside the
// clamped source rectangle reading transparent. The four samples are
// interpolated in premultiplied space so an opaque pixel beside a transparent
// one keeps its colour instead of rounding to a dark edge.
func bilinearAt(src *Canvas, x0, y0, x1, y1 int, s, t float64) Color {
	ix := int(math.Floor(s - 0.5))
	iy := int(math.Floor(t - 0.5))
	fx := s - 0.5 - float64(ix)
	fy := t - 0.5 - float64(iy)

	var pr, pg, pb, pa float64
	for j, wy := range [2]float64{1 - fy, fy} {
		for i, wx := range [2]float64{1 - fx, fx} {
			xx, yy := ix+i, iy+j
			if xx < x0 || yy < y0 || xx >= x1 || yy >= y1 {
				continue
			}
			w := wx * wy
			c := src.At(xx, yy)
			a := float64(c.A())
			pm := a / 255 * w
			pr += float64(c.R()) * pm
			pg += float64(c.G()) * pm
			pb += float64(c.B()) * pm
			pa += a * w
		}
	}
	if pa < 0.5 {
		return 0
	}
	to := func(v float64) uint8 {
		v = v * 255 / pa
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v + 0.5)
	}
	return RGBA(to(pr), to(pg), to(pb), uint8(pa+0.5))
}
