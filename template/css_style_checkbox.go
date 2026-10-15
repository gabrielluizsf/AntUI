package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Checkbox(p CheckboxParams) {
	st := cs.styleFor(css.RoleCheckbox, p.State)
	if p.On {
		st = cs.styleFor(css.RoleCheckbox, State{On: true, Focused: p.State.Focused, Hovered: p.State.Hovered})
	}
	u := cs.u()
	box := boxSize * u
	cs.paintBox(p.Win, st, p.X, p.Y, box, box)
	if p.On {
		mark := cs.ink(st)
		cv := p.Win.Canvas()
		cv.Line(canvas.Point{X: float64(float64(float64(p.X+4*u))), Y: float64(float64(float64(p.Y+9*u)))}, canvas.Point{X: float64(float64(float64(p.X+7*u))), Y: float64(float64(float64(p.Y+13*u)))}, mark)
		cv.Line(canvas.Point{X: float64(float64(float64(p.X+5*u))), Y: float64(float64(float64(p.Y+9*u)))}, canvas.Point{X: float64(float64(float64(p.X+8*u))), Y: float64(float64(float64(p.Y+13*u)))}, mark)
		cv.Line(canvas.Point{X: float64(float64(float64(p.X+7*u))), Y: float64(float64(float64(p.Y+13*u)))}, canvas.Point{X: float64(float64(float64(p.X+14*u))), Y: float64(float64(float64(p.Y+5*u)))}, mark)
		cv.Line(canvas.Point{X: float64(float64(float64(p.X+8*u))), Y: float64(float64(float64(p.Y+13*u)))}, canvas.Point{X: float64(float64(float64(p.X+15*u))), Y: float64(float64(float64(p.Y+5*u)))}, mark)
	}
	cs.drawText(p.Win, st, p.X+box+labelGap*u, p.Y, box, box, p.Label)
}
