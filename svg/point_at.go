package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// pointAt is where a distance along an outline lands: the point itself and the
// way the outline is going there, which is the direction the letter at that
// place faces. Past the end of the outline the last point and the last
// direction are what the writing hangs on, since an outline shorter than the
// writing on it still has to put that writing somewhere, and a distance before
// the start runs back along the first stretch for the same reason.
func pointAt(p *canvas.Path, d float64) (x, y, angle float64) {
	pts, closed := p.Points()
	for i, sub := range pts {
		segs := len(sub) - 1
		if closed[i] {
			segs++
		}
		for j := 0; j < segs; j++ {
			a, b := sub[j], sub[j+1]
			if j == segs-1 && closed[i] {
				b = sub[0]
			}
			seg := math.Hypot(b.X-a.X, b.Y-a.Y)
			if seg <= 0 {
				continue
			}
			if d <= seg {
				t := d / seg
				return a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t,
					math.Atan2(b.Y-a.Y, b.X-a.X)
			}
			d -= seg
			x, y, angle = b.X, b.Y, math.Atan2(b.Y-a.Y, b.X-a.X)
		}
	}
	return x, y, angle
}
