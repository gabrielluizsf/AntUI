package template

import "github.com/gabrielluizsf/antui/canvas"

// vivid is what the theme draws the ink in, sharpened for the window's
// opacity: text and the neon details get lifted toward white as the window
// fades, so a translucent window — whose whole frame the compositor dims —
// keeps its fonts and its button details bright and readable. The dark panels
// are left alone: the contrast between ink and panel is exactly what the lift
// is saving. At full opacity it is the palette as it sits.
type vivid struct {
	text, neon, dim, edge canvas.Color
	glow                  float64 // how much the halos are amplified, ≥ 1
}
