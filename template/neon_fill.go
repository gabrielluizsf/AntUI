package template

import "github.com/gabrielluizsf/antui/canvas"

// neonFill paints a control's body and its border: a dark panel with a thin
// coloured edge. What happens under the pointer — who glows, who brightens —
// is the caller's call, because a button and a field answer differently.
func neonFill(cv *canvas.Canvas, x, y, w, h, radius int, body, border canvas.Color) {
	cv.FillRoundRect(x, y, w, h, radius, body)
	cv.RoundRect(x, y, w, h, radius, border)
}
