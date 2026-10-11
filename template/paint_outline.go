package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintOutline rings the box outside its border, pushed out by outline-offset.
// The ring is drawn side by side, so its corners are square rather than mitered.
func (cs *cssStyle) paintOutline(win *antui.Window, st css.Style, x, y, w, h int) {
	u := cs.u()
	off := st.OutlineOffset * u
	ow := st.OutlineWidth * u
	col := cs.effectColor(win, st, st.OutlineColor)
	ox, oy := x-off-ow, y-off-ow
	ow2, oh2 := w+2*(off+ow), h+2*(off+ow)
	for side := 0; side < 4; side++ {
		cs.paintSide(win.Canvas(), side, st.OutlineStyle, ow, ox, oy, ow2, oh2, col)
	}
}
