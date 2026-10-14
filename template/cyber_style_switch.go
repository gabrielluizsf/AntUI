package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Switch(p SwitchParams) {
	win := p.Win
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	trackW, trackH := switchTrackW*u, switchTrackH*u

	if p.State.Focused || p.State.Hovered {
		glow(cv, halo{x: p.X, y: p.Y, w: trackW, h: trackH, radius: trackH, colour: v.neon, strength: int(float64(30) * v.glow), u: u})
	}
	body := cyberPalette.panel
	if p.State.Hovered {
		body = cyberPalette.panelHov
	}
	border := v.edge
	if p.On || p.State.Focused {
		border = v.neon
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, trackH/2), body)
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, trackH/2), border)

	if p.On {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackH, Height: trackH}, trackH/2), canvas.Fade(v.neon, 120))
	}
	colour := v.edge
	if p.On {
		colour = v.neon
	}
	cx := p.X + trackH/2
	if p.On {
		cx = p.X + trackW - trackH/2
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: p.Y + trackH/2, R: trackH/2 - 3*u}, colour)
	cv.Circle(canvas.Circle{X: cx, Y: p.Y + trackH/2, R: trackH/2 - 4*u}, cyberPalette.panelAct)
	cv.TextScaled(p.X+trackW+labelGap*u, p.Y+(trackH-textHeight(u))/2, p.Label, v.text, u)
}
