package canvas

import "github.com/gabrielluizsf/antui/font"

// TextStyle is how a run of text is drawn beyond its colour: an integer
// enlargement of the face, a fake bold (the stroke drawn twice), a fake
// italic (each glyph row shifted to lean right) and the extra spacing CSS
// asks for after every glyph and after every space.
type TextStyle struct {
	Scale         int // 0 or 1 draw the face at its own size
	Bold          bool
	Italic        bool
	LetterSpacing int // extra pixels after every glyph
	WordSpacing   int // extra pixels after every space
}

func (o TextStyle) norm() TextStyle {
	if o.Scale < 1 {
		o.Scale = 1
	}
	return o
}

// TextWidthStyled is [TextWidth] under a [TextStyle], measured with the
// program's default face.
func TextWidthStyled(text string, o TextStyle) int {
	return faceWidth(DefaultFace(), text, o)
}

// TextWidthStyled is how wide a string is under a style in this canvas's face.
func (cv *Canvas) TextWidthStyled(text string, o TextStyle) int {
	return faceWidth(cv.face(), text, o)
}

// faceWidth sums one glyph at a time so the per-glyph and per-space spacing
// land exactly where Draw puts them.
func faceWidth(f *Face, text string, o TextStyle) int {
	o = o.norm()
	fs := f.Scaled(o.Scale)
	widest, line := 0, 0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, line)
			line = 0
		case '\r':
		case '\t':
			line += fs.Width("\t")
		case ' ':
			line += fs.Width(" ") + o.WordSpacing
		default:
			line += fs.Width(string(r)) + o.LetterSpacing
		}
	}
	return max(widest, line)
}

// DrawStyled writes a string with a [TextStyle], the baseline of the first
// line at y, and returns the width of the widest line. It is [Face.Draw] with
// the weight, slant and spacing the stylesheet asked for.
func (cv *Canvas) DrawStyled(x, y int, text string, c Color, o TextStyle) int {
	o = o.norm()
	f := cv.face()
	if f != BuiltinFace() && f.file != nil {
		sf := f.Scaled(o.Scale)
		return cv.drawStyledMask(sf, x, y+sf.Ascent(), text, c, o)
	}
	return cv.drawStyledBitmap(x, y, text, c, o)
}

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
		m := f.glyph(r)
		if m.W > 0 {
			blitMaskStyled(cv, m, int(penX+0.5), int(penY+0.5), c, o)
		}
		penX += f.advance(r) + float64(o.LetterSpacing)
	}
	return max(widest, int(penX-float64(x)+0.5))
}

func (cv *Canvas) strokeCell(x, y, scale int, c Color, bold bool) {
	if scale <= 1 {
		cv.Pixel(x, y, c)
		if bold {
			cv.Pixel(x+1, y, c)
		}
		return
	}
	cv.FillRect(x, y, scale, scale, c)
	if bold {
		cv.FillRect(x+scale, y, scale, scale, c)
	}
}

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
			cv.Pixel(px, top+row, RGBA(c.R(), c.G(), c.B(), uint8(v*alpha/255)))
			if o.Bold {
				cv.Pixel(px+step, top+row, RGBA(c.R(), c.G(), c.B(), uint8(v*alpha/255)))
			}
		}
	}
}
