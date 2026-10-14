package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Triangle draws the edges of a triangle.
func (win *Window) Triangle(points TrianglePoints, c canvas.Color) {
	p0, p1, p2 := points.P0, points.P1, points.P2
	win.cv.Triangle(int(p0.X), int(p0.Y), int(p1.X), int(p1.Y), int(p2.X), int(p2.Y), c)
}
