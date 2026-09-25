package canvas

import "math"

// RingArc draws a thick arc of the annulus between radius-thickness and radius,
// centred on cx, cy. The arc starts at start radians and sweeps sweep radians;
// a positive sweep turns clockwise, because the canvas y axis grows downward.
// Both the radial and the angular edges are antialiased. A sweep of at least a
// full turn paints the whole ring.
func (cv *Canvas) RingArc(cx, cy, radius, thickness int, start, sweep float64, c Color) {
	if radius <= 0 || thickness <= 0 || sweep == 0 {
		return
	}
	thickness = min(thickness, radius)
	inner := radius - thickness
	base := int(c.A())
	if base == 0 {
		return
	}
	sweep = max(min(sweep, 2*math.Pi), -2*math.Pi)
	dir, span := 1.0, sweep
	if sweep < 0 {
		dir, span = -1, -sweep
	}
	full := span >= 2*math.Pi
	feather := 1.0 / float64(max(radius, 1))

	x0 := max(cx-radius-1, cv.Clip.X)
	y0 := max(cy-radius-1, cv.Clip.Y)
	x1 := min(cx+radius+1, cv.Clip.X+cv.Clip.Width-1)
	y1 := min(cy+radius+1, cv.Clip.Y+cv.Clip.Height-1)
	for y := y0; y <= y1; y++ {
		dy := y - cy
		for x := x0; x <= x1; x++ {
			dx := x - cx
			radial := ellipseCoverage(dx, dy, radius, radius)
			if radial == 0 {
				continue
			}
			if inner > 0 {
				radial -= ellipseCoverage(dx, dy, inner, inner)
				if radial <= 0 {
					continue
				}
			}
			acov := 1.0
			if !full {
				delta := math.Mod(dir*(math.Atan2(float64(dy), float64(dx))-start), 2*math.Pi)
				if delta < 0 {
					delta += 2 * math.Pi
				}
				if delta > span {
					continue
				}
				if delta < feather {
					acov = delta / feather
				}
				if span-delta < feather {
					acov = min(acov, (span-delta)/feather)
				}
				if acov <= 0 {
					continue
				}
			}
			cv.Pixel(x, y, Fade(c, int(float64(radial)*acov)*base/255))
		}
	}
}
