package canvas

func (o TextStyle) norm() TextStyle {
	if o.Scale < 1 {
		o.Scale = 1
	}
	return o
}
