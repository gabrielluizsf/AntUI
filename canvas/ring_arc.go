package canvas

import "math"

// RingArcParams holds parameters for drawing a thick ring arc.
type RingArcParams struct {
	Center     Point
	Radius     int
	Thickness  int
	Start      float64
	Sweep      float64
	Color      Color
}

// RingArc draws a thick arc of the annulus.
func (cv *Canvas) RingArc(p RingArcParams) {
	cx, cy := int(p.Center.X), int(p.Center.Y)
	radius, thickness := p.Radius, p.Thickness
	start, sweep := p.Start, p.Sweep
	c := p.Color
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
			acov := arcSweepCoverage(dx, dy, start, span, dir, feather, full)
			if acov <= 0 {
				continue
			}
			cv.Pixel(x, y, Fade(c, int(float64(radial)*acov)*base/255))
		}
	}
}