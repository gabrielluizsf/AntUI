package canvas

// drawStyledBitmap is the built-in font with slant and a doubled stroke.
func (cv *Canvas) drawStyledBitmap(x, y int, text string, c Color, o TextStyle) int {
	scale := o.Scale
	penX, penY, widest := x, y, 0

	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, penX-x)
			penX, penY = x, penY+FontHeight*scale
			continue
		case '\r':
			continue
		case '\t':
			penX += FontWidth * scale * 4
			continue
		}
		if r == ' ' {
			penX += FontWidth*scale + o.WordSpacing
			continue
		}
		rows := glyph(r)
		for row := range FontHeight {
			bits := rows[row]
			if bits == 0 {
				continue
			}
			shift := 0
			if o.Italic {
				shift = (FontHeight - 1 - row) * scale / 4
			}
			for col := range FontWidth {
				if bits&(0x80>>col) == 0 {
					continue
				}
				dx := penX + col*scale + shift
				cv.strokeCell(dx, penY+row*scale, scale, c, o.Bold)
			}
		}
		penX += FontWidth*scale + o.LetterSpacing
	}
	return max(widest, penX-x)
}
