package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Switch(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	trackW, trackH := switchTrackW*u, switchTrackH*u

	if s.Focused || s.Hovered {
		glow(cv, halo{x: x, y: y, w: trackW, h: trackH, radius: trackH, colour: v.neon, strength: int(float64(30) * v.glow), u: u})
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	border := v.edge
	if on || s.Focused {
		border = v.neon
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: trackW, Height: trackH}, trackH/2), body)
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: trackW, Height: trackH}, trackH/2), border)

	if on {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: trackH, Height: trackH}, trackH/2), canvas.Fade(v.neon, 120))
	}
	colour := v.edge
	if on {
		colour = v.neon
	}
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: y + trackH/2, R: trackH/2 - 3*u}, colour)
	cv.Circle(canvas.Circle{X: cx, Y: y + trackH/2, R: trackH/2 - 4*u}, cyberPalette.panelAct)
	cv.TextScaled(x+trackW+labelGap*u, y+(trackH-textHeight(u))/2, label, v.text, u)
}
