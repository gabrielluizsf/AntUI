package canvas

// Draw writes a string into a canvas with the baseline of the first line at
// y, and returns the width of the widest line.
func (f *Face) Draw(cv *Canvas, x, y int, text string, c Color) int {
	if f == nil || f.file == nil {
		return cv.textBitmap(x, y-f.Ascent(), text, c, f.scaleOf())
	}
	if cv == nil {
		return f.Width(text)
	}
	penX, penY := float64(x), y
	widest := 0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, int(penX-float64(x)+0.5))
			penX, penY = float64(x), penY+f.height
			continue
		case '\r':
			continue
		case '\t':
			penX += f.advance(' ') * 4
			continue
		case ' ':
			penX += f.advance(' ')
			continue
		}
		if cg := f.colorGlyph(r); cg != nil {
			f.blitColor(cv, cg, penX, float64(penY), TextStyle{})
		} else if m := f.glyph(r); m.W > 0 {
			blitMask(cv, m, int(penX+0.5), penY, c)
		}
		penX += f.advance(r)
	}
	return max(widest, int(penX-float64(x)+0.5))
}
