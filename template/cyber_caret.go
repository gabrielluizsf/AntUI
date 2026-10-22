package template

import "github.com/gabrielluizsf/antui/canvas"

// caretGlyph is where the site's text cursor sits, how tall it stands and
// what colour it burns, at the scale the field around it draws at.
type caretGlyph struct {
	X, Y, H int
	U       int
	Color   canvas.Color
}

// drawCyberCaret paints the site's own text cursor: a glowing blade rather
// than the hairline every other field draws, with a white-hot tip the way a
// welding arc or a lightsaber reads brightest right at the point.
func drawCyberCaret(cv *canvas.Canvas, g caretGlyph) {
	w := 2 * g.U
	cv.BoxShadow(canvas.BoxShadowParams{
		X: g.X - g.U, Y: g.Y - g.U, Width: w + 2*g.U, Height: g.H + 2*g.U,
		Blur: 3 * g.U, Color: canvas.Fade(g.Color, 140),
	})
	cv.FillRect(canvas.Area{X: g.X, Y: g.Y, Width: w, Height: g.H}, g.Color)
	cv.FillRect(canvas.Area{X: g.X, Y: g.Y, Width: w, Height: max(g.H/3, g.U)}, canvas.White)
}
