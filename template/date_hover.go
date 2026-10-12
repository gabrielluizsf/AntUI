package template

import "github.com/gabrielluizsf/antui"

// dateHover carries the highlight to the day the pointer is over, but only on
// frames the pointer actually moved, so the arrow keys keep control while the
// mouse sits still.
func (t *uiTemplate) dateHover(win *antui.Window, v dateView, g dateGeom) dateView {
	if win.Hovered(g.px, g.py, g.dw, g.dh) && win.MouseMoved() {
		if d := dayAt(g.px, g.py, g.u, win.MouseX(), win.MouseY(), v.year, v.month); d > 0 {
			v.day = d
		}
	}
	return v
}
