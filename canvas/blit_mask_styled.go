package canvas

import "github.com/gabrielluizsf/antui/font"

func blitMaskStyled(cv *Canvas, m *font.Mask, x, baseline int, c Color, o TextStyle) {
	alpha := int(c.A())
	if alpha == 0 {
		return
	}
	top := baseline - m.Top
	step := max(o.Scale, 1)
	for row := range m.H {
		shift := 0
		if o.Italic {
			shift = (m.H - 1 - row) / 4
		}
		for col := range m.W {
			v := int(m.Alpha[row*m.Stride+col])
			if v == 0 {
				continue
			}
			px := x + m.Left + col + shift
			cv.Pixel(px, top+row, RGBA(RGBAParams{R: c.R(), G: c.G(), B: c.B(), A: uint8(v * alpha / 255)}))
			if o.Bold {
				cv.Pixel(px+step, top+row, RGBA(RGBAParams{R: c.R(), G: c.G(), B: c.B(), A: uint8(v * alpha / 255)}))
			}
		}
	}
}
