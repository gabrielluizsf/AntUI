package canvas

import "github.com/gabrielluizsf/antui/font"

// blitMask composites one glyph's coverage onto a canvas in a colour.
func blitMask(cv *Canvas, m *font.Mask, x, baseline int, c Color) {
	alpha := int(c.A())
	if alpha == 0 {
		return
	}
	top := baseline - m.Top
	for row := range m.H {
		for col := range m.W {
			v := int(m.Alpha[row*m.Stride+col])
			if v == 0 {
				continue
			}
			cv.Pixel(x+m.Left+col, top+row, RGBA(RGBAParams{R: c.R(), G: c.G(), B: c.B(), A: uint8(v * alpha / 255)}))
		}
	}
}
