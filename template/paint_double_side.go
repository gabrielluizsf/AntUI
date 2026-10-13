package template

import "github.com/gabrielluizsf/antui/canvas"

// paintDoubleSide draws a double border: two thin lines with a gap between.
// A border too thin for a real double falls back to a solid side.
func (cs *cssStyle) paintDoubleSide(cv *canvas.Canvas, s boxSide) {
	line := max(s.width/3, 1)
	if s.width-2*line < 1 {
		cs.paintSolidSide(cv, s)
		return
	}
	// The outer line is flush with the box's edge; the inner one is a stroke-
	// width in from the opposite edge, and the gap between the two is what
	// makes the pair read as a double rather than as a thick border. The two
	// offsets are the same on all four sides, because `off` is measured from
	// whichever edge this side is flush with.
	for _, off := range [2]int{0, s.width - line} {
		x, y, w, h := s.rect(0, s.along(), off, line)
		cv.FillRect(x, y, w, h, s.colour)
	}
}
