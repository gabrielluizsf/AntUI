package template

import "github.com/gabrielluizsf/antui/canvas"

func (builtinStyle) Label(p LabelParams) {
	p.Win.Canvas().TextScaled(canvas.Point{X: float64(p.X), Y: float64(p.Y)}, p.Text, canvas.TextStyle{Color: p.Win.Theme().Text, Scale: Scale(p.Win)})
}
