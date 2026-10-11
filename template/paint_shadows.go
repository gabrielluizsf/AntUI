package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintShadows paints the box's shadows on the given side of the surface. A
// drop-shadow filter is painted with the outer shadows, since both fall
// behind the box.
func (cs *cssStyle) paintShadows(win *antui.Window, st css.Style, x, y, w, h, rx, ry int, inset bool) {
	if !inset {
		for i := range st.Filters {
			if st.Filters[i].Drop != nil {
				cs.paintShadow(win, st, *st.Filters[i].Drop, x, y, w, h, rx, ry, false)
			}
		}
	}
	for i := range st.BoxShadow {
		if st.BoxShadow[i].Inset == inset {
			cs.paintShadow(win, st, st.BoxShadow[i], x, y, w, h, rx, ry, inset)
		}
	}
}
