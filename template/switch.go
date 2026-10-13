package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Switch(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	trackW, trackH := switchTrackW*u, switchTrackH*u

	fill := th.Surface
	if s.Hovered {
		fill = th.SurfaceHover
	}
	border := th.Border
	if on || s.Focused {
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: trackW, Height: trackH}, trackH/2), fill)
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: trackW, Height: trackH}, trackH/2), border)

	colour := th.Border
	if on {
		colour = th.Accent
	}
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: y + trackH/2, R: trackH/2 - 2*u}, colour)
	if s.Focused {
		cv.Circle(canvas.Circle{X: x + trackW/2, Y: y + trackH/2, R: trackH/2 + 3*u}, th.Accent)
	}
	cv.TextScaled(x+trackW+labelGap*u, y+(trackH-textHeight(u))/2, label, th.Text, u)
}
