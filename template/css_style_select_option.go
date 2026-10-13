package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	st := cs.styleFor(css.RoleSelect, s)
	if selected {
		st = cs.styleFor(css.RoleSelect, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	fill := cs.fill(st)
	if selected || s.Hovered {
		fill = canvas.Shade(fill, hoverFillShade)
	}
	u := cs.u()
	cv := win.Canvas()
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: h}, cs.fade(st, fill))
	if selected || s.Hovered {
		cv.FillRect(canvas.Area{X: x, Y: y, Width: 3 * u, Height: h}, cs.ink(st))
	}
	cs.drawText(win, st, x, y, w, h, label)
}
