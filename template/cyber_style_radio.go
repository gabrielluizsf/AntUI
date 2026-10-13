package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	size := boxSize * u
	cx, cy := x+size/2, y+size/2
	if s.Focused || s.Hovered {
		glow(cv, halo{x: x, y: y, w: size, h: size, radius: size / 2, colour: v.neon, strength: int(float64(30) * v.glow), u: u})
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: size / 2}, body)
	border := v.edge
	if on || s.Hovered {
		border = v.neon
	}
	cv.Circle(canvas.Circle{X: cx, Y: cy, R: size / 2}, border)
	if on {
		cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: size/2 - 4*u}, v.neon)
	}
	cv.TextScaled(x+size+labelGap*u, y+(size-textHeight(u))/2, label, v.text, u)
}
