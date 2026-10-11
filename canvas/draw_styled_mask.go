package canvas

// drawStyledMask is the TrueType path: the glyph mask is blitted row by row,
// leaning and doubled the same way.
func (cv *Canvas) drawStyledMask(f *Face, x, y int, text string, c Color, o TextStyle) int {
	penX, penY := float64(x), float64(y)
	widest := 0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, int(penX-float64(x)+0.5))
			penX, penY = float64(x), penY+float64(f.height)
			continue
		case '\r':
			continue
		case '\t':
			penX += f.advance(' ') * 4
			continue
		case ' ':
			penX += f.advance(' ') + float64(o.WordSpacing)
			continue
		}
		if cg := f.colorGlyph(r); cg != nil {
			f.blitColor(cv, cg, penX, penY, o)
		} else if m := f.glyph(r); m.W > 0 {
			blitMaskStyled(cv, m, int(penX+0.5), int(penY+0.5), c, o)
		}
		penX += f.advance(r) + float64(o.LetterSpacing)
	}
	return max(widest, int(penX-float64(x)+0.5))
}
