package template

import "github.com/gabrielluizsf/antui/canvas"

func (builtinStyle) Label(p LabelParams) {
	cv := p.Win.Canvas()
	// A label this long or this scaled up can measure wider than the box it
	// was given; clipping to that box is what keeps it from painting over
	// whatever sits next to it rather than just running past its own edge.
	savedClip := cv.Clip
	if p.W > 0 && p.H > 0 {
		cv.SetClip(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H})
	}
	cv.TextScaled(canvas.Point{X: float64(p.X), Y: float64(p.Y)}, p.Text, canvas.TextStyle{Color: p.Win.Theme().Text, Scale: Scale(p.Win)})
	cv.Clip = savedClip
}
