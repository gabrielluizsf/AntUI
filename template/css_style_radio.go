package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Radio(p RadioParams) {
	st := cs.styleFor(css.RoleRadio, p.State)
	if p.On {
		st = cs.styleFor(css.RoleRadio, State{On: true, Focused: p.State.Focused, Hovered: p.State.Hovered})
	}
	u := cs.u()
	size := boxSize * u
	half := size / 2
	cv := p.Win.Canvas()
	cv.FillCircle(canvas.Circle{X: p.X + half, Y: p.Y + half, R: half}, cs.fill(st))
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
	cv.Circle(canvas.Circle{X: p.X + half, Y: p.Y + half, R: half}, cs.fade(st, border))
	if p.On {
		cv.FillCircle(canvas.Circle{X: p.X + half, Y: p.Y + half, R: half - 4*u}, cs.ink(st))
	}
	cs.drawText(p.Win, st, p.X+size+labelGap*u, p.Y, size, size, p.Label)
}
