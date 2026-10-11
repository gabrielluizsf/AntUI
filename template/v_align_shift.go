package template

import "github.com/gabrielluizsf/antui/template/css"

// vAlignShift moves the text top by vertical-align within its box.
func (cs *cssStyle) vAlignShift(st css.Style, th, h, padT, padB int) int {
	center := (h - padT - padB - th) / 2
	switch st.VerticalAlign {
	case css.VerticalAlignSub:
		return th / 4
	case css.VerticalAlignSuper:
		return -th / 4
	case css.VerticalAlignMiddle:
		return -th / 8
	case css.VerticalAlignTop, css.VerticalAlignTextTop:
		return -center
	case css.VerticalAlignBottom, css.VerticalAlignTextBottom:
		return center
	}
	if !st.BaselineShift.Auto() && !st.BaselineShift.None() {
		return -st.BaselineShift.Resolve(cs.units(st)) * cs.u()
	}
	return 0
}
