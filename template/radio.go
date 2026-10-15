package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Radio(p RadioParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	size := boxSize * u

	surface := th.Surface
	if p.State.Hovered {
		surface = th.SurfaceHover
	}
	border := th.Border
	if p.On {
		border = th.Accent
	}
	cv.FillCircle(canvas.Circle{X: p.X + size/2, Y: p.Y + size/2, R: size / 2}, surface)
	cv.Circle(canvas.Circle{X: p.X + size/2, Y: p.Y + size/2, R: size / 2}, border)
	if p.State.Focused {
		cv.Circle(canvas.Circle{X: p.X + size/2, Y: p.Y + size/2, R: size/2 + 3*u}, th.Accent)
	}
	if p.On {
		drawIconIn(cv, IconDot, p.X+4*u, p.Y+4*u, 10*u, th.Accent)
	}
	cv.TextScaled(canvas.Point{X: float64(p.X + size + labelGap*u), Y: float64(p.Y + (size-textHeight(u))/2)}, p.Label, canvas.TextStyle{Color: th.Text, Scale: u})
}
