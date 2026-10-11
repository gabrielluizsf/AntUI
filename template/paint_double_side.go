package template

import "github.com/gabrielluizsf/antui/canvas"

// paintDoubleSide draws a double border: two thin lines with a gap between.
// A border too thin for a real double falls back to a solid side.
func (cs *cssStyle) paintDoubleSide(cv *canvas.Canvas, side int, bw, x, y, w, h int, col canvas.Color) {
	line := max(bw/3, 1)
	if bw-2*line < 1 {
		cs.paintSolidSide(cv, side, bw, x, y, w, h, col)
		return
	}
	switch side {
	case 0:
		cv.FillRect(x, y, w, line, col)
		cv.FillRect(x, y+bw-line, w, line, col)
	case 1:
		cv.FillRect(x+w-line, y, line, h, col)
		cv.FillRect(x+w-bw, y, line, h, col)
	case 2:
		cv.FillRect(x, y+h-line, w, line, col)
		cv.FillRect(x, y+h-bw, w, line, col)
	default:
		cv.FillRect(x, y, line, h, col)
		cv.FillRect(x+bw-line, y, line, h, col)
	}
}
