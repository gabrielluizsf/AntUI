package svg

import "github.com/gabrielluizsf/antui/canvas"

// movedBox is a rectangle taken through a transform: the four corners go and
// the box around where they land is what comes back, which is exact for the
// transforms that keep a rectangle a rectangle and never smaller than the
// truth for the ones that turn it into a diamond.
func movedBox(x, y, w, h float64, m canvas.Matrix) (float64, float64, float64, float64) {
	x0, y0 := m.Map(x, y)
	x1, y1 := m.Map(x+w, y)
	x2, y2 := m.Map(x, y+h)
	x3, y3 := m.Map(x+w, y+h)
	minX, maxX := min(x0, x1, x2, x3), max(x0, x1, x2, x3)
	minY, maxY := min(y0, y1, y2, y3), max(y0, y1, y2, y3)
	return minX, minY, maxX - minX, maxY - minY
}
