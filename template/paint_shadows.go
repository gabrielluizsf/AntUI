package template

import "github.com/gabrielluizsf/antui"

// paintShadows paints the box's shadows on the given side of the surface. A
// drop-shadow filter is painted with the outer shadows, since both fall
// behind the box.
func (cs *cssStyle) paintShadows(win *antui.Window, b shadowBox) {
	if !b.inset {
		for i := range b.st.Filters {
			if b.st.Filters[i].Drop != nil {
				cs.paintShadow(win, *b.st.Filters[i].Drop, b)
			}
		}
	}
	for i := range b.st.BoxShadow {
		if b.st.BoxShadow[i].Inset == b.inset {
			cs.paintShadow(win, b.st.BoxShadow[i], b)
		}
	}
}
