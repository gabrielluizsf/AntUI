package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Progress(p ProgressParams) {
	win := p.Win
	cv := win.Canvas()
	v := vividFor(win)

	value := min(max(p.Value, 0), 1)
	filled := int(float32(p.W)*value + 0.5)

	// The site's progress bar: a near-invisible rail and a flat accent fill,
	// square corners, no glow.
	cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, canvas.Fade(canvas.White, 13))
	if filled > 0 {
		cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: filled, Height: p.H}, v.neon)
	}
}
