package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintOutline rings the box outside its border, pushed out by outline-offset.
// Solid rings follow the box radius; a zero-radius box stays square.
func (cs *cssStyle) paintOutline(win *antui.Window, st css.Style, x, y, w, h int) {
	u := cs.u()
	off := st.OutlineOffset * u
	ow := st.OutlineWidth * u
	col := cs.effectColor(win, st, st.OutlineColor)
	rx := st.Radius[css.SideTop] * u
	ry := rx
	if st.RadiusY[css.SideTop] != 0 {
		ry = st.RadiusY[css.SideTop] * u
	}
	ox, oy := x-off-ow, y-off-ow
	ow2, oh2 := w+2*(off+ow), h+2*(off+ow)
	if st.OutlineStyle == css.BorderSolid {
		// A solid outline is one rounded ring. Painting four rectangular sides
		// separately makes their corners overlap and compounds a translucent
		// focus colour, which is the dark square halo visible around inputs.
		for i := range ow {
			rw, rh := ow2-2*i, oh2-2*i
			if rw <= 0 || rh <= 0 {
				break
			}
			orx, ory := 0, 0
			if rx > 0 {
				orx = max(rx-i, 0)
			}
			if ry > 0 {
				ory = max(ry-i, 0)
			}
			cv := win.Canvas()
			cv.RoundRectXY(ox+i, oy+i, rw, rh, orx, ory, col)
		}
		return
	}
	for side := 0; side < 4; side++ {
		cs.paintSide(win.Canvas(), boxSide{side: side, x: ox, y: oy, w: ow2, h: oh2,
			width: ow, colour: col}, st.OutlineStyle)
	}
}
