package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberMarks draws the graduations along the trough, lit up to the knob.
func paintCyberMarks(cv *canvas.Canvas, v vivid, trackX, trackY, trackW, trackH, knobX, u int) {
	for i := 1; i < 4; i++ {
		gx := trackX + trackW*i/4

		tickColor := canvas.Fade(v.dim, 130)

		if gx <= knobX {
			tickColor = v.neon
		}

		cv.FillRoundRect(gx-u, trackY-3*u, 2*u, trackH+6*u, u, tickColor)
	}
}
