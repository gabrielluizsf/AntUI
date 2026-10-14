package template

func (builtinStyle) Label(p LabelParams) {
	p.Win.Canvas().TextScaled(p.X, p.Y, p.Text, p.Win.Theme().Text, Scale(p.Win))
}
