package canvas

// Height is how far apart two lines are under this style, in pixels: the
// height of the face it draws in, enlarged by the style's scale the same way
// Draw puts the lines down. A style that named no face keeps the height the
// default face has always had.
func (o TextStyle) Height() int {
	o = o.norm()
	if o.Face == nil {
		return TextHeight() * o.Scale
	}
	return o.Face.Scaled(o.Scale).Height()
}
