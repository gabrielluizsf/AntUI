package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Triangle draws the edges of a triangle.
func (win *Window) Triangle(points TrianglePoints, c canvas.Color) {
	p0, p1, p2 := points.P0, points.P1, points.P2
	win.cv.Triangle([3]canvas.Point{{X: float64(int(p0.X)), Y: float64(int(p0.Y))}, {X: float64(int(p1.X)), Y: float64(int(p1.Y))}, {X: float64(int(p2.X)), Y: float64(int(p2.Y))}}, c)
}
