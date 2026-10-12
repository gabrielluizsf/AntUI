package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberTrack draws the trough: a faint glow, a dark glass body with a
// thin edge, and the lit fill up to the knob.
func paintCyberTrack(cv *canvas.Canvas, v vivid, trackX, trackY, trackW, trackH, knobX, u int) {
	// Very subtle glow behind the entire track.
	glow(cv, trackX-2*u, trackY-2*u, trackW+4*u, trackH+4*u, trackH/2+2*u, v.neon, 5, u)

	// Dark glass body.
	cv.FillRoundRect(trackX, trackY, trackW, trackH, trackH/2, canvas.Fade(cyberPalette.panelAct, 230))

	// Thin outer edge.
	cv.RoundRect(trackX, trackY, trackW, trackH, trackH/2, v.edge)

	fillW := knobX - trackX
	if fillW > 0 {
		// Glow around the active portion.
		glow(cv, trackX-u, trackY-u, fillW+2*u, trackH+2*u, trackH/2+u, v.neon, 10, u)
		cv.FillRoundRect(trackX, trackY, fillW, trackH, trackH/2, v.neon)
	}
}
