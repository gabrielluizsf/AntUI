package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// transformPad is how far, in device pixels, a widget's pixels can reach past
// its border box: shadows and outlines. Whatever is closer to the box is
// inside the region already.
func (cs *cssStyle) transformPad(st css.Style) int {
	u := cs.u()
	pad := 0
	span := func(s css.Shadow) int {
		ex := max(s.X, -s.X)
		if s.Y > ex {
			ex = s.Y
		}
		if -s.Y > ex {
			ex = -s.Y
		}
		return ex + s.Spread + s.Blur + 2
	}
	for i := range st.Filters {
		if f := st.Filters[i]; f.Drop != nil {
			pad = max(pad, span(*f.Drop))
		} else if f.Kind == canvas.FilterBlur {
			pad = max(pad, int(f.Amount)+1)
		}
	}
	for i := range st.BoxShadow {
		if s := st.BoxShadow[i]; !s.Inset {
			pad = max(pad, span(s))
		}
	}
	for i := range st.TextShadow {
		pad = max(pad, span(st.TextShadow[i]))
	}
	if cs.outlineOn(st) {
		pad = max(pad, st.OutlineOffset+st.OutlineWidth)
	}
	if pad < 0 {
		return 0
	}
	return pad * u
}
