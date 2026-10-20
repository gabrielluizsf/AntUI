package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// paintCyberTrack draws the trough the way the site draws progress: a dark
// rail with a hairline frame and a flat accent fill, no glow.
func paintCyberTrack(cv *canvas.Canvas, v vivid, trackX, trackY, trackW, trackH, knobX int) {
	area := canvas.Round(canvas.Area{X: trackX, Y: trackY, Width: trackW, Height: trackH}, trackH/2)
	cv.FillRoundRect(area, cyberPalette.field)
	cv.RoundRect(area, v.edge)
	fillW := knobX - trackX
	if fillW > 0 {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: trackX, Y: trackY, Width: fillW, Height: trackH}, trackH/2), v.neon)
	}
}
