package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintShadow draws one shadow, scaling its lengths with the window.
func (cs *cssStyle) paintShadow(win *antui.Window, sh css.Shadow, b shadowBox) {
	u := cs.u()
	col := cs.effectColor(win, b.st, sh.Color)
	win.Canvas().BoxShadow(b.x, b.y, b.w, b.h, b.rx, b.ry,
		sh.X*u, sh.Y*u, sh.Blur*u, sh.Spread*u, col, b.inset)
}
