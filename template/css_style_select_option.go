package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) SelectOption(p SelectOptionParams) {
	st := cs.styleFor(css.RoleSelect, p.State)
	if p.Selected {
		st = cs.styleFor(css.RoleSelect, State{On: true, Focused: p.State.Focused, Hovered: p.State.Hovered})
	}
	fill := cs.fill(st)
	if p.Selected || p.State.Hovered {
		fill = canvas.Shade(fill, hoverFillShade)
	}
	u := cs.u()
	cv := p.Win.Canvas()
	cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, cs.fade(st, fill))
	if p.Selected || p.State.Hovered {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: 3 * u, Height: p.H}, cs.ink(st))
	}
	cs.drawText(p.Win, st, p.X, p.Y, p.W, p.H, p.Label)
}
