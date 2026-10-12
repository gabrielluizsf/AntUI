package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintSliderTrack draws the groove and the filled part up to the thumb.
func paintSliderTrack(cv *canvas.Canvas, th *antui.Theme, trackX, trackY, trackW, trackH, knobX int) {
	// Unfilled part.
	cv.FillRoundRect(trackX, trackY, trackW, trackH, trackH/2, th.Border)

	// Filled part.
	fillW := knobX - trackX
	if fillW > 0 {
		cv.FillRoundRect(trackX, trackY, fillW, trackH, trackH/2, th.Accent)
	}
}
