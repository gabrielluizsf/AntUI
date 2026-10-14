package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// FillTriangle fills a triangle.
func (win *Window) FillTriangle(points TrianglePoints, c canvas.Color) {
	p0, p1, p2 := points.P0, points.P1, points.P2
	win.cv.FillTriangle(int(p0.X), int(p0.Y), int(p1.X), int(p1.Y), int(p2.X), int(p2.Y), c)
}
