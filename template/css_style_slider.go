package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Slider(p SliderParams) {
	st := cs.styleFor(css.RoleSlider, p.State)
	u := cs.u()
	value := clampFloat(p.Value, 0, 1)
	knobRadius := max(p.H/2, 6*u)
	usable := max(p.W-2*knobRadius, 1)
	knobX := p.X + knobRadius + int(value*float32(usable)+0.5)
	trackY := p.Y + p.H/2 - 3*u

	cv := p.Win.Canvas()
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X + knobRadius, Y: trackY, Width: knobX - (p.X + knobRadius), Height: 6 * u}, 3*u), cs.fill(st))
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X + knobRadius, Y: trackY, Width: usable, Height: 6 * u}, 3*u), canvas.Shade(cs.fill(st), trackShade))
	cv.FillCircle(canvas.Circle{X: knobX, Y: p.Y + p.H/2, R: knobRadius}, cs.ink(st))
	cv.FillCircle(canvas.Circle{X: knobX, Y: p.Y + p.H/2, R: knobRadius - 4*u}, cs.fill(st))
}
