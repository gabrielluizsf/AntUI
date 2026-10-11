package canvas

import "math"

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
