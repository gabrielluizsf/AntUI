package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Switch(p SwitchParams) {
	st := cs.styleFor(css.RoleSwitch, p.State)
	if p.On {
		st = cs.styleFor(css.RoleSwitch, State{On: true, Focused: p.State.Focused, Hovered: p.State.Hovered})
	}
	u := cs.u()
	trackW, trackH := switchTrackW*u, switchTrackH*u

	fill := cs.fill(st)
	if p.On {
		fill = cs.ink(st)
	}
	cv := p.Win.Canvas()
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, trackH/2), cs.fade(st, fill))
	cx := p.X + trackH/2
	if p.On {
		cx = p.X + trackW - trackH/2
	}
	thumb := canvas.Shade(fill, thumbShade)
	if p.On {
		thumb = cs.fill(st)
	}
	if thumb.A() == 0 {
		// The thumb is the widget's, not the sheet's background: a sheet that
		// gave the track no colour of its own still gets a switch that is off,
		// drawn as the dot at the end of an empty track.
		thumb = cs.ink(st)
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: p.Y + trackH/2, R: trackH/2 - 2*u}, thumb)
	cs.run(p.Win.Canvas(), st, p.X+trackW+labelGap*u, p.Y+(trackH-cs.textHeight(st))/2, p.Label)
}
