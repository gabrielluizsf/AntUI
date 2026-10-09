package css

import "github.com/gabrielluizsf/antui/canvas"

// boxStyle carries the box-model and border fields of a Style: how big the
// box is, how far it is padded, how thick and what colour its border is.
type boxStyle struct {
	// Box model.
	Width, Height        Length
	MinWidth, MaxWidth   Length
	MinHeight, MaxHeight Length
	Margin               [4]Length // top, right, bottom, left
	Padding              [4]Length
	BoxSizing            bool // border-box when true, content-box otherwise

	// Border. BoxStyle[i] is one of the Border* constants.
	BorderWidth [4]int
	BoxStyle    [4]uint8
	BoxColor    [4]canvas.Color
	Radius      [4]int
	RadiusY     [4]int // vertical radii; zero means the same as Radius
}

// BorderOn reports whether any side paints a border.
func (s *Style) BorderOn() bool {
	return s.BoxStyle[0] != BorderNone || s.BoxStyle[1] != BorderNone ||
		s.BoxStyle[2] != BorderNone || s.BoxStyle[3] != BorderNone
}
