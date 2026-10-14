package template

func (cyberStyle) Label(p LabelParams) {
	p.Win.Canvas().TextScaled(p.X, p.Y, p.Text, vividFor(p.Win).text, Scale(p.Win))
}
