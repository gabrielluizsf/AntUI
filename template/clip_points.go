package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// clipPoints resolves a clip-path polygon against a box: each coordinate
// expression is evaluated with the box's width on the x axis and its height on
// the y axis, so "50%" and "calc(100% - 6px)" both land where the site's own
// clip-path would. It returns nil for a box with no polygon, or one with fewer
// than three resolvable points.
func clipPoints(st css.Style, x, y, w, h int) []canvas.Point {
	if len(st.ClipPath) < 3 {
		return nil
	}
	xu := css.Units{Width: w, Height: h, Font: h, Root: h}
	yu := css.Units{Width: h, Height: h, Font: h, Root: h}
	pts := make([]canvas.Point, 0, len(st.ClipPath))
	for _, p := range st.ClipPath {
		px, ok := css.EvalMath(p[0], xu)
		if !ok {
			return nil
		}
		py, ok := css.EvalMath(p[1], yu)
		if !ok {
			return nil
		}
		pts = append(pts, canvas.Point{X: float64(x) + px, Y: float64(y) + py})
	}
	return pts
}
