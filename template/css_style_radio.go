package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	st := cs.styleFor(css.RoleRadio, s)
	if on {
		st = cs.styleFor(css.RoleRadio, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	u := cs.u()
	size := boxSize * u
	half := size / 2
	cv := win.Canvas()
	cv.FillCircle(canvas.Circle{X: x + half, Y: y + half, R: half}, cs.fill(st))
	border := st.BoxColor[css.SideTop]
	if !st.BorderOn() {
		border = canvas.Shade(cs.fill(st), hoverBorderShade)
	}
	if border.A() == 0 {
		// The ring is the widget's: a sheet that gave the radio neither a border
		// nor a background still gets one it can be picked by, in the ink it
		// labels itself with.
		border = cs.ink(st)
	}
	cv.Circle(canvas.Circle{X: x + half, Y: y + half, R: half}, cs.fade(st, border))
	if on {
		cv.FillCircle(canvas.Circle{X: x + half, Y: y + half, R: half - 4*u}, cs.ink(st))
	}
	cs.drawText(win, st, x+size+labelGap*u, y, size, size, label)
}
