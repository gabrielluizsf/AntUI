package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintShadow draws one shadow, scaling its lengths with the window.
func (cs *cssStyle) paintShadow(win *antui.Window, st css.Style, sh css.Shadow, x, y, w, h, rx, ry int, inset bool) {
	u := cs.u()
	col := cs.effectColor(win, st, sh.Color)
	win.Canvas().BoxShadow(x, y, w, h, rx, ry, sh.X*u, sh.Y*u, sh.Blur*u, sh.Spread*u, col, inset)
}
