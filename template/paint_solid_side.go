package template

import "github.com/gabrielluizsf/antui/canvas"

// paintSolidSide fills one straight side.
func (cs *cssStyle) paintSolidSide(cv *canvas.Canvas, side, bw, x, y, w, h int, col canvas.Color) {
	switch side {
	case 0:
		cv.FillRect(x, y, w, bw, col)
	case 1:
		cv.FillRect(x+w-bw, y, bw, h, col)
	case 2:
		cv.FillRect(x, y+h-bw, w, bw, col)
	default:
		cv.FillRect(x, y, bw, h, col)
	}
}
