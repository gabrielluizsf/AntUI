package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintShadow draws one shadow, scaling its lengths with the window.
func (cs *cssStyle) paintShadow(win *antui.Window, sh css.Shadow, b shadowBox) {
	u := cs.u()
	col := cs.effectColor(win, b.st, sh.Color)
	win.Canvas().BoxShadow(canvas.BoxShadowParams{
		X: b.x, Y: b.y, Width: b.w, Height: b.h, RX: b.rx, RY: b.ry,
		DX: sh.X * u, DY: sh.Y * u, Blur: sh.Blur * u, Spread: sh.Spread * u,
		Color: col, Inset: b.inset,
	})
}
