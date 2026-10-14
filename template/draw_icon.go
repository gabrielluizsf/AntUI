package template

import "github.com/gabrielluizsf/antui/canvas"

// DrawIcon paints a named icon into the rectangle given, in the colour given,
// and says whether there was one to paint.
//
// The icon is fitted into the rectangle keeping its own proportions and centred
// in what is left over, so a square rectangle gives an icon drawn square and a
// rectangle that is not square gives the same icon with its margins, which is
// what every icon set does and what keeps a tick from being stretched into a
// different mark.
//
// Painting nothing is not an error: an icon that is not registered simply does
// not appear, so a widget that asks for one it may not have still draws.
func DrawIcon(cv *canvas.Canvas, name string, x, y, w, h int, tint canvas.Color) bool {
	if cv == nil {
		return false
	}
	painted := IconCanvas(name, w, h, tint)
	if painted == nil {
		return false
	}
	cv.BlitScaled(canvas.Area{X: x, Y: y, Width: w, Height: h}, painted)
	return true
}
