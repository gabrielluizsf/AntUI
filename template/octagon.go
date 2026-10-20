package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// octagon is the site's panel window: a rectangle with all four corners cut
// by cut pixels, the shape its inputs and dialogs wear.
func octagon(x, y, w, h, cut int) []canvas.Point {
	cut = min(cut, min(w, h)/2)
	return []canvas.Point{
		{X: float64(x + cut), Y: float64(y)},
		{X: float64(x + w - cut), Y: float64(y)},
		{X: float64(x + w), Y: float64(y + cut)},
		{X: float64(x + w), Y: float64(y + h - cut)},
		{X: float64(x + w - cut), Y: float64(y + h)},
		{X: float64(x + cut), Y: float64(y + h)},
		{X: float64(x), Y: float64(y + h - cut)},
		{X: float64(x), Y: float64(y + cut)},
	}
}

// octagonFill paints a panel body and its one-unit edge as the site's
// chamfered window.
func octagonFill(cv *canvas.Canvas, x, y, w, h, cut, u int, body, border canvas.Color) {
	pts := octagon(x, y, w, h, cut)
	fillClip(cv, pts, body)
	strokeClip(cv, pts, u, border)
}
