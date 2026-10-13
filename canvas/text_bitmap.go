package canvas

import "math"

// textBitmap is the built-in font at a size: each cell of its shape is drawn
// as a square as large as that size says, and the cells that come out no
// pixels wide at a size smaller than its own are left out rather than pulled
// open to one, so writing thinner than the shape thins instead of growing.
func (cv *Canvas) textBitmap(x, y int, text string, c Color, unit float64) int {
	if unit <= 0 {
		unit = 1
	}
	penX, penY := float64(x), float64(y)
	widest := 0.0

	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, penX-float64(x))
			penX, penY = float64(x), penY+float64(FontHeight)*unit
			continue
		case '\t':
			penX += float64(FontWidth) * unit * 4
			continue
		case '\r':
			continue
		}
		if r != ' ' {
			rows := glyph(r)
			for row := range FontHeight {
				bits := rows[row]
				if bits == 0 {
					continue
				}
				top := int(math.Round(penY + float64(row)*unit))
				bottom := int(math.Round(penY + float64(row+1)*unit))
				if bottom <= top {
					continue
				}
				for col := range FontWidth {
					if bits&(0x80>>col) == 0 {
						continue
					}
					left := int(math.Round(penX + float64(col)*unit))
					right := int(math.Round(penX + float64(col+1)*unit))
					if right <= left {
						continue
					}
					if right-left == 1 && bottom-top == 1 {
						cv.Pixel(left, top, c)
					} else {
						cv.FillRect(Area{left, top, right - left, bottom - top}, c)
					}
				}
			}
		}
		penX += float64(FontWidth) * unit
	}
	return int(math.Round(max(widest, penX-float64(x))))
}
