package template

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// strokeClip draws a polygon's edge as a one-unit ring: each edge becomes a
// quad between the outline and a copy shrunk toward the centroid.
func strokeClip(cv *canvas.Canvas, pts []canvas.Point, u int, c canvas.Color) {
	if u <= 0 || len(pts) < 3 {
		return
	}
	cx, cy := 0.0, 0.0
	for _, p := range pts {
		cx += p.X
		cy += p.Y
	}
	n := float64(len(pts))
	cx, cy = cx/n, cy/n

	inner := make([]canvas.Point, len(pts))
	for i, p := range pts {
		dx, dy := cx-p.X, cy-p.Y
		d := math.Sqrt(dx*dx + dy*dy)
		k := 1.0
		if d > 0 {
			k = max(0, d-float64(u)) / d
		}
		inner[i] = canvas.Point{X: cx - dx*k, Y: cy - dy*k}
	}
	for i := range pts {
		j := (i + 1) % len(pts)
		cv.FillTriangle([3]canvas.Point{pts[i], pts[j], inner[j]}, c)
		cv.FillTriangle([3]canvas.Point{pts[i], inner[j], inner[i]}, c)
	}
}
