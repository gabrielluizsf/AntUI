package template

import "github.com/gabrielluizsf/antui/canvas"

// DrawIconIn is [DrawIcon] at the size a widget draws its own marks at, centred
// in the box given. It is what the built-in widgets call, so that the mark they
// draw is the same size and in the same place whatever the scale.
func drawIconIn(cv *canvas.Canvas, name string, x, y, size int, tint canvas.Color) bool {
	return drawIcon(cv, name, x, y, size, size, tint)
}
