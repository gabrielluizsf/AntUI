package template

import "github.com/gabrielluizsf/antui"

func (builtinStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	size := boxSize * u

	surface := th.Surface
	if s.Hovered {
		surface = th.SurfaceHover
	}
	border := th.Border
	if on {
		border = th.Accent
	}
	cv.FillCircle(x+size/2, y+size/2, size/2, surface)
	cv.Circle(x+size/2, y+size/2, size/2, border)
	if s.Focused {
		cv.Circle(x+size/2, y+size/2, size/2+3*u, th.Accent)
	}
	if on {
		DrawIconIn(cv, IconDot, x+4*u, y+4*u, 10*u, th.Accent)
	}
	cv.TextScaled(x+size+labelGap*u, y+(size-textHeight(u))/2, label, th.Text, u)
}
