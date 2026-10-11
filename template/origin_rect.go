package template

import "github.com/gabrielluizsf/antui/template/css"

// originRect insets the box by the border for a padding box, or the border
// plus padding for a content box; the border box is the whole box. Each
// background attaches to its own positioning area.
func (cs *cssStyle) originRect(st css.Style, x, y, w, h int, origin uint8) (int, int, int, int) {
	u := cs.u()
	b := func(side int) int { return st.BorderWidth[side] * u }
	p := func(side int) int { return cs.base(st, st.Padding[side]) }
	switch origin {
	case css.BackBorder:
		return x, y, w, h
	case css.BackContent:
		return x + b(3) + p(3), y + b(0) + p(0),
			max(w-b(1)-b(3)-p(1)-p(3), 0), max(h-b(0)-b(2)-p(0)-p(2), 0)
	default: // the padding box is the initial background-origin
		return x + b(3), y + b(0), max(w-b(1)-b(3), 0), max(h-b(0)-b(2), 0)
	}
}
