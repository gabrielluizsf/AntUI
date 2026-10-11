package template

import "github.com/gabrielluizsf/antui/canvas"

// unionArea is the bounding box of two areas; an empty one is absorbed.
func unionArea(a, b canvas.Area) canvas.Area {
	if b.Width <= 0 || b.Height <= 0 {
		return a
	}
	if a.Width <= 0 || a.Height <= 0 {
		return b
	}
	x1 := max(a.X+a.Width, b.X+b.Width)
	y1 := max(a.Y+a.Height, b.Y+b.Height)
	a.X = min(a.X, b.X)
	a.Y = min(a.Y, b.Y)
	a.Width, a.Height = x1-a.X, y1-a.Y
	return a
}
