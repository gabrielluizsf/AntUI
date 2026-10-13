package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Switch(win *antui.Window, s State, x, y int, label string, on bool) {
	st := cs.styleFor(css.RoleSwitch, s)
	if on {
		st = cs.styleFor(css.RoleSwitch, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	u := cs.u()
	trackW, trackH := switchTrackW*u, switchTrackH*u

	fill := cs.fill(st)
	if on {
		fill = cs.ink(st)
	}
	cv := win.Canvas()
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: trackW, Height: trackH}, trackH/2), cs.fade(st, fill))
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	thumb := canvas.Shade(fill, thumbShade)
	if on {
		thumb = cs.fill(st)
	}
	if thumb.A() == 0 {
		// The thumb is the widget's, not the sheet's background: a sheet that
		// gave the track no colour of its own still gets a switch that is off,
		// drawn as the dot at the end of an empty track.
		thumb = cs.ink(st)
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: y + trackH/2, R: trackH/2 - 2*u}, thumb)
	cs.run(cv, st, x+trackW+labelGap*u, y+(trackH-cs.textHeight(st))/2, label)
}
