package template

import "github.com/gabrielluizsf/antui/canvas"

func (builtinStyle) Switch(p SwitchParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	trackW, trackH := switchTrackW*u, switchTrackH*u

	fill := th.Surface
	if p.State.Hovered {
		fill = th.SurfaceHover
	}
	border := th.Border
	if p.On || p.State.Focused {
		border = th.Accent
	}
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, trackH/2), fill)
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: trackW, Height: trackH}, trackH/2), border)

	colour := th.Border
	if p.On {
		colour = th.Accent
	}
	cx := p.X + trackH/2
	if p.On {
		cx = p.X + trackW - trackH/2
	}
	cv.FillCircle(canvas.Circle{X: cx, Y: p.Y + trackH/2, R: trackH/2 - 2*u}, colour)
	if p.State.Focused {
		cv.Circle(canvas.Circle{X: p.X + trackW/2, Y: p.Y + trackH/2, R: trackH/2 + 3*u}, th.Accent)
	}
	cv.TextScaled(p.X+trackW+labelGap*u, p.Y+(trackH-textHeight(u))/2, p.Label, th.Text, u)
}
