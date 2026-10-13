package template

import "github.com/gabrielluizsf/antui/canvas"

// paintDashSide draws a dashed or dotted side: dashes as rectangles, dots as
// small discs, spaced along the side. The two differ only in what is laid down
// and how far apart, so they share the walk and the spacing.
func (cs *cssStyle) paintDashSide(cv *canvas.Canvas, s boxSide, dotted bool) {
	length := s.along()
	u := cs.u()
	step := max(3*s.width, 3*u)
	seg := max(2*s.width, 2*u)
	if dotted {
		step = max(2*s.width, 2*u)
		seg = s.width
	}

	// A dot is a disc of the border's own thickness, so it has to fit whole
	// before it is drawn and one hanging off the end is dropped; a dash is a
	// rectangle clipped to whatever is left of the side.
	for at := 0; at < length; at += step {
		if dotted && at+s.width > length {
			break
		}
		n := min(seg, length-at)
		if n <= 0 {
			break
		}
		x, y, w, h := s.rect(at, n, 0, s.width)
		if dotted {
			cv.FillCircle(x+w/2, y+h/2, max(s.width/2, 1), s.colour)
		} else {
			cv.FillRect(x, y, w, h, s.colour)
		}
	}
}
