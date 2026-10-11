package template

import "github.com/gabrielluizsf/antui/canvas"

// paintDashSide draws a dashed or dotted side: dashes as rectangles, dots as
// small discs, spaced along the side.
func (cs *cssStyle) paintDashSide(cv *canvas.Canvas, side int, bw, x, y, w, h int, col canvas.Color, dotted bool) {
	horizontal := side == 0 || side == 2
	length := w
	if !horizontal {
		length = h
	}
	u := cs.u()
	step := max(3*bw, 3*u)
	if dotted {
		step = max(2*bw, 2*u)
	}

	if dotted {
		r := max(bw/2, 1)
		for p := 0; p+bw <= length; p += step {
			if horizontal {
				yy := y
				if side == 2 {
					yy = y + h - bw
				}
				cv.FillCircle(x+p+bw/2, yy+bw/2, r, col)
			} else {
				xx := x
				if side == 1 {
					xx = x + w - bw
				}
				cv.FillCircle(xx+bw/2, y+p+bw/2, r, col)
			}
		}
		return
	}

	seg := max(2*bw, 2*u)
	for p := 0; p < length; p += step {
		n := min(seg, length-p)
		if n <= 0 {
			break
		}
		if horizontal {
			yy := y
			if side == 2 {
				yy = y + h - bw
			}
			cv.FillRect(x+p, yy, n, bw, col)
		} else {
			xx := x
			if side == 1 {
				xx = x + w - bw
			}
			cv.FillRect(xx, y+p, bw, n, col)
		}
	}
}
