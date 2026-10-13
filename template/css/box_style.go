package css

import "github.com/gabrielluizsf/antui/canvas"

// boxStyle carries the box-model and border fields of a Style: how big the
// box is, how far it is padded, how thick and what colour its border is.
type boxStyle struct {
	// Box model.
	Width, Height        Length
	MinWidth, MaxWidth   Length
	MinHeight, MaxHeight Length
	Margin               [sideCount]Length // top, right, bottom, left
	Padding              [sideCount]Length
	BoxSizing            bool // border-box when true, content-box otherwise

	// Border. BoxStyle[i] is one of the Border* constants, indexed by side.
	BorderWidth [sideCount]int
	BoxStyle    [sideCount]uint8
	BoxColor    [sideCount]canvas.Color
	Radius      [sideCount]int
	RadiusY     [sideCount]int // vertical radii; zero means the same as Radius
}

// BorderOn reports whether any side paints a border.
func (s *Style) BorderOn() bool {
	return s.BoxStyle[SideTop] != BorderNone || s.BoxStyle[SideRight] != BorderNone ||
		s.BoxStyle[SideBottom] != BorderNone || s.BoxStyle[SideLeft] != BorderNone
}
