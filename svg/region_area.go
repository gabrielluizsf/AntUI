package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// regionArea is where a mask reaches on the canvas: the four corners of the
// region go through the transform that puts the drawing on it, and the box
// around where they land is what is cut to — exact for the transforms that
// keep a rectangle a rectangle, and never smaller than the truth for the ones
// that do not, which is what a region has to be.
func regionArea(x, y, w, h float64, m canvas.Matrix) (x0, y0, x1, y1 int) {
	ax, ay := m.Map(x, y)
	bx, by := m.Map(x+w, y)
	cx, cy := m.Map(x, y+h)
	dx, dy := m.Map(x+w, y+h)
	minX := math.Floor(min(ax, bx, cx, dx))
	minY := math.Floor(min(ay, by, cy, dy))
	maxX := math.Ceil(max(ax, bx, cx, dx))
	maxY := math.Ceil(max(ay, by, cy, dy))
	return int(minX), int(minY), int(maxX), int(maxY)
}
