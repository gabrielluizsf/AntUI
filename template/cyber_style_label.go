package template

import "github.com/gabrielluizsf/antui/canvas"

func (cyberStyle) Label(p LabelParams) {
	p.Win.Canvas().TextScaled(canvas.Point{X: float64(p.X), Y: float64(p.Y)}, p.Text, canvas.TextStyle{Color: vividFor(p.Win).text, Scale: Scale(p.Win)})
}
