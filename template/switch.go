package template

import "github.com/gabrielluizsf/antui"

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
	cv.FillRoundRect(x, y, trackW, trackH, trackH/2, fill)
	cv.RoundRect(x, y, trackW, trackH, trackH/2, border)

	colour := th.Border
	if on {
		colour = th.Accent
	}
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	cv.FillCircle(cx, y+trackH/2, trackH/2-2*u, colour)
	if s.Focused {
		cv.Circle(x+trackW/2, y+trackH/2, trackH/2+3*u, th.Accent)
	}
	cv.TextScaled(x+trackW+labelGap*u, y+(trackH-textHeight(u))/2, label, th.Text, u)
}
