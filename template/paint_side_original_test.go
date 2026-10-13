package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// TestBorderPaintersMatchTheOriginalsBitForBit is the safety net for the
// boxSide rewrite: every style, on every side, at several widths and several
// box shapes, must land on exactly the same pixels the original painters did.
func TestBorderPaintersMatchTheOriginalsBitForBit(t *testing.T) {
	const pad = 6
	boxes := []struct{ x, y, w, h int }{
		{pad, pad, 60, 30},
		{pad, pad, 31, 47}, // odd, and not a multiple of anything
		{pad, pad, 9, 7},   // shorter than one dash on some widths
		{pad, pad, 120, 90},
	}
	widths := []int{1, 2, 3, 4, 6, 9}
	styles := []uint8{css.BorderSolid, css.BorderDashed, css.BorderDotted, css.BorderDouble}

	for _, b := range boxes {
		for _, bw := range widths {
			for _, style := range styles {
				for _, side := range []int{css.SideTop, css.SideRight, css.SideBottom, css.SideLeft} {
					name := styleName(style) + "/" + sideName(side)
					s := boxSide{side: side, x: b.x, y: b.y, w: b.w, h: b.h,
						width: bw, colour: canvas.Black}

					winNew, cvNew, err := antui.Offscreen(b.x+b.w+pad*2, b.y+b.h+pad*2)
					if err != nil {
						t.Fatal(err)
					}
					cs := &cssStyle{win: winNew}
					winNew.Begin()
					winNew.Clear(canvas.White)
					cs.paintSide(cvNew, s, style)
					winNew.End()

					winOld, cvOld, err := antui.Offscreen(b.x+b.w+pad*2, b.y+b.h+pad*2)
					if err != nil {
						t.Fatal(err)
					}
					winOld.Begin()
					winOld.Clear(canvas.White)
					switch style {
					case css.BorderDouble:
						origDouble(cvOld, side, bw, b.x, b.y, b.w, b.h, canvas.Black)
					case css.BorderDashed:
						origDash(cvOld, side, bw, b.x, b.y, b.w, b.h, canvas.Black, false, 1)
					case css.BorderDotted:
						origDash(cvOld, side, bw, b.x, b.y, b.w, b.h, canvas.Black, true, 1)
					default:
						origSolid(cvOld, side, bw, b.x, b.y, b.w, b.h, canvas.Black)
					}
					winOld.End()

					if a, c := snapshot(cvNew), snapshot(cvOld); a != c {
						t.Errorf("box %dx%d width %d %s: painted differently than before",
							b.w, b.h, bw, name)
					}
				}
			}
		}
	}
}

func styleName(s uint8) string {
	switch s {
	case css.BorderDashed:
		return "dashed"
	case css.BorderDotted:
		return "dotted"
	case css.BorderDouble:
		return "double"
	default:
		return "solid"
	}
}

func sideName(s int) string {
	switch s {
	case css.SideRight:
		return "right"
	case css.SideBottom:
		return "bottom"
	case css.SideLeft:
		return "left"
	default:
		return "top"
	}
}

// The four border painters were rewritten to take a boxSide instead of eight
// loose numbers, and the dashed and dotted cases were folded onto one walk. The
// bodies below are the originals, copied verbatim from before the rewrite, kept
// here so the new ones can be held to them: a border painted a pixel off is
// invisible to every other test in the package.

// origSolid is the original paintSolidSide, untouched.
func origSolid(cv *canvas.Canvas, side, bw, x, y, w, h int, col canvas.Color) {
	switch side {
	case 0:
		cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: bw}, col)
	case 1:
		cv.FillRect(canvas.Area{X: x + w - bw, Y: y, Width: bw, Height: h}, col)
	case 2:
		cv.FillRect(canvas.Area{X: x, Y: y + h - bw, Width: w, Height: bw}, col)
	default:
		cv.FillRect(canvas.Area{X: x, Y: y, Width: bw, Height: h}, col)
	}
}

// origDouble is the original paintDoubleSide, untouched.
func origDouble(cv *canvas.Canvas, side, bw, x, y, w, h int, col canvas.Color) {
	line := max(bw/3, 1)
	if bw-2*line < 1 {
		origSolid(cv, side, bw, x, y, w, h, col)
		return
	}
	switch side {
	case 0:
		cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: line}, col)
		cv.FillRect(canvas.Area{X: x, Y: y + bw - line, Width: w, Height: line}, col)
	case 1:
		cv.FillRect(canvas.Area{X: x + w - line, Y: y, Width: line, Height: h}, col)
		cv.FillRect(canvas.Area{X: x + w - bw, Y: y, Width: line, Height: h}, col)
	case 2:
		cv.FillRect(canvas.Area{X: x, Y: y + h - line, Width: w, Height: line}, col)
		cv.FillRect(canvas.Area{X: x, Y: y + h - bw, Width: w, Height: line}, col)
	default:
		cv.FillRect(canvas.Area{X: x, Y: y, Width: line, Height: h}, col)
		cv.FillRect(canvas.Area{X: x + bw - line, Y: y, Width: line, Height: h}, col)
	}
}

// origDash is the original paintDashSide, untouched.
func origDash(cv *canvas.Canvas, side, bw, x, y, w, h int, col canvas.Color, dotted bool, u int) {
	horizontal := side == 0 || side == 2
	length := w
	if !horizontal {
		length = h
	}
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
				cv.FillCircle(canvas.Circle{X: x + p + bw/2, Y: yy + bw/2, R: r}, col)
			} else {
				xx := x
				if side == 1 {
					xx = x + w - bw
				}
				cv.FillCircle(canvas.Circle{X: xx + bw/2, Y: y + p + bw/2, R: r}, col)
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
			cv.FillRect(canvas.Area{X: x + p, Y: yy, Width: n, Height: bw}, col)
		} else {
			xx := x
			if side == 1 {
				xx = x + w - bw
			}
			cv.FillRect(canvas.Area{X: xx, Y: y + p, Width: bw, Height: n}, col)
		}
	}
}

// snapshot is every pixel of a canvas, as a comparable string.
func snapshot(cv *canvas.Canvas) string {
	b := make([]byte, 0, cv.Width*cv.Height*4)
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			c := cv.At(x, y)
			b = append(b, byte(c.R()), byte(c.G()), byte(c.B()), byte(c.A()))
		}
	}
	return string(b)
}
