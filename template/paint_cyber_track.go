package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberTrack draws the trough: a faint glow, a dark glass body with a
// thin edge, and the lit fill up to the knob.
func paintCyberTrack(cv *canvas.Canvas, v vivid, trackX, trackY, trackW, trackH, knobX, u int) {
	// Very subtle glow behind the entire track.
	glow(cv, halo{x: trackX - 2*u, y: trackY - 2*u, w: trackW + 4*u, h: trackH + 4*u, radius: trackH/2 + 2*u, colour: v.neon, strength: 5, u: u})

	// Dark glass body.
	cv.FillRoundRect(trackX, trackY, trackW, trackH, trackH/2, canvas.Fade(cyberPalette.panelAct, 230))

	// Thin outer edge.
	cv.RoundRect(trackX, trackY, trackW, trackH, trackH/2, v.edge)

	fillW := knobX - trackX
	if fillW > 0 {
		// Glow around the active portion.
		glow(cv, halo{x: trackX - u, y: trackY - u, w: fillW + 2*u, h: trackH + 2*u, radius: trackH/2 + u, colour: v.neon, strength: 10, u: u})
		cv.FillRoundRect(trackX, trackY, fillW, trackH, trackH/2, v.neon)
	}
}
