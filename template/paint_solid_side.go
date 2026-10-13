package template

import "github.com/gabrielluizsf/antui/canvas"

// paintSolidSide fills one straight side.
func (cs *cssStyle) paintSolidSide(cv *canvas.Canvas, s boxSide) {
	x, y, w, h := s.rect(0, s.along(), 0, s.width)
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: h}, s.colour)
}
