package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Slider(win *antui.Window, s State, x, y, w, h int, value float32) {
	st := cs.styleFor(css.RoleSlider, s)
	u := cs.u()
	value = clampFloat(value, 0, 1)
	knobRadius := max(h/2, 6*u)
	usable := max(w-2*knobRadius, 1)
	knobX := x + knobRadius + int(value*float32(usable)+0.5)
	trackY := y + h/2 - 3*u

	cv := win.Canvas()
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x + knobRadius, Y: trackY, Width: knobX - (x + knobRadius), Height: 6 * u}, 3*u), cs.fill(st))
	cv.FillRoundRect(canvas.Round(canvas.Area{X: x + knobRadius, Y: trackY, Width: usable, Height: 6 * u}, 3*u), canvas.Shade(cs.fill(st), trackShade))
	cv.FillCircle(canvas.Circle{X: knobX, Y: y + h/2, R: knobRadius}, cs.ink(st))
	cv.FillCircle(canvas.Circle{X: knobX, Y: y + h/2, R: knobRadius - 4*u}, cs.fill(st))
}
