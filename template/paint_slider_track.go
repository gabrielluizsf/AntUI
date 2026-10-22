package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// sliderTrack is the groove's own geometry and where the thumb sits along it.
type sliderTrack struct {
	X, Y, Width, Height int
	KnobX               int
}

// paintSliderTrack draws the groove and the filled part up to the thumb.
func paintSliderTrack(cv *canvas.Canvas, th *antui.Theme, t sliderTrack) {
	radius := t.Height / 2
	// Unfilled part.
	cv.FillRoundRect(canvas.Round(canvas.Area{X: t.X, Y: t.Y, Width: t.Width, Height: t.Height}, radius), th.SurfaceActive)

	// Filled part.
	fillW := t.KnobX - t.X
	if fillW > 0 {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: t.X, Y: t.Y, Width: fillW, Height: t.Height}, radius), th.Accent)
	}
}
