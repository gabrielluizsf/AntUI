package canvas

import "github.com/gabrielluizsf/antui/font"

// FontWidth and FontHeight are the cell size of the built-in font, in pixels.
var (
	FontWidth  = font.Width  // 8
	FontHeight = font.Height // 16
)

// glyph returns the 16 rows of the built-in font for a codepoint.
func glyph(r rune) []byte { return font.Glyph(r) }

// Put writes a colour with no blending, which is what a rasterizer that has
// already decided what it wants needs. It honours the clip.
func (cv *Canvas) Put(x, y int, c Color) {
	if !cv.insideClip(x, y) {
		return
	}
	if cv.Pixels != nil {
		cv.Pixels[y*cv.Stride+x] = c
		return
	}
	if cv.narrow != nil {
		cv.narrowWrite(x, y, c)
	}
}

// Pixel blends a colour over what is already there, honouring the clip.
func (cv *Canvas) Pixel(x, y int, c Color) {
	if !cv.insideClip(x, y) {
		return
	}
	if cv.Pixels == nil {
		if cv.narrow == nil {
			return
		}
		switch c.A() {
		case 0:
			return
		case 255:
			cv.narrowWrite(x, y, c)
		default:
			cv.narrowWrite(x, y, Blend(cv.narrowRead(x, y), c))
		}
		return
	}
	at := y*cv.Stride + x
	cv.Pixels[at] = Blend(cv.Pixels[at], c)
}

// At reads a pixel, in any format. Points outside the canvas read as zero.
// Unlike drawing, reading ignores the clip.
func (cv *Canvas) At(x, y int) Color {
	if x < 0 || y < 0 || x >= cv.Width || y >= cv.Height {
		return 0
	}
	if cv.Pixels == nil {
		if cv.narrow == nil {
			return 0
		}
		return cv.narrowRead(x, y)
	}
	return cv.Pixels[y*cv.Stride+x]
}

// FillRect fills a rectangle, blending when the colour is translucent.
func (cv *Canvas) FillRect(x, y, w, h int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	x0 := max(x, cv.Clip.X)
	y0 := max(y, cv.Clip.Y)
	x1 := min(x+w, cv.Clip.X+cv.Clip.Width)
	y1 := min(y+h, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}

	if cv.Pixels == nil {
		if cv.narrow == nil {
			return
		}
		switch c.A() {
		case 0:
		case 255:
			for iy := y0; iy < y1; iy++ {
				cv.narrowRun(x0, x1, iy, c)
			}
		default:
			for iy := y0; iy < y1; iy++ {
				for ix := x0; ix < x1; ix++ {
					cv.Pixel(ix, iy, c)
				}
			}
		}
		return
	}

	switch c.A() {
	case 0:
	case 255:
		for iy := y0; iy < y1; iy++ {
			row := cv.Pixels[iy*cv.Stride : iy*cv.Stride+cv.Width]
			for ix := x0; ix < x1; ix++ {
				row[ix] = c
			}
		}
	default:
		for iy := y0; iy < y1; iy++ {
			BlendRowSolid(cv.Pixels[iy*cv.Stride+x0:iy*cv.Stride+x1], c)
		}
	}
}

// Clear paints the whole clip region. The colour is forced opaque.
func (cv *Canvas) Clear(c Color) {
	cv.FillRect(cv.Clip.X, cv.Clip.Y, cv.Clip.Width, cv.Clip.Height, c|0xFF000000)
}

// Rect draws a one-pixel outline.
func (cv *Canvas) Rect(x, y, w, h int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	cv.FillRect(x, y, w, 1, c)
	if h > 1 {
		cv.FillRect(x, y+h-1, w, 1, c)
	}
	if h > 2 {
		cv.FillRect(x, y+1, 1, h-2, c)
		if w > 1 {
			cv.FillRect(x+w-1, y+1, 1, h-2, c)
		}
	}
}

// Line draws a line with Bresenham's algorithm.
func (cv *Canvas) Line(x0, y0, x1, y1 int, c Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy
	for {
		cv.Pixel(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// discCoverage is the antialiased coverage of a disc, 0 outside to 255 in.
func discCoverage(dx, dy, radius int) int {
	dist2 := uint64(dx*dx + dy*dy)
	dist := int64(isqrt64(dist2 << 16))
	edge := int64(radius)*256 + 128
	cover := edge - dist
	if cover <= 0 {
		return 0
	}
	if cover >= 256 {
		return 255
	}
	return int(cover)
}

// FillCircle draws a filled, antialiased disc.
func (cv *Canvas) FillCircle(cx, cy, radius int, c Color) {
	if radius <= 0 {
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}
	x0 := max(cx-radius-1, cv.Clip.X)
	y0 := max(cy-radius-1, cv.Clip.Y)
	x1 := min(cx+radius+1, cv.Clip.X+cv.Clip.Width-1)
	y1 := min(cy+radius+1, cv.Clip.Y+cv.Clip.Height-1)

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			alpha := discCoverage(x-cx, y-cy, radius)
			if alpha == 0 {
				continue
			}
			cv.Pixel(x, y, Fade(c, alpha*baseAlpha/255))
		}
	}
}

// Circle draws a one-pixel circle outline.
func (cv *Canvas) Circle(cx, cy, radius int, c Color) {
	if radius <= 0 {
		return
	}
	x, y := radius, 0
	err := 1 - radius
	for x >= y {
		cv.Pixel(cx+x, cy+y, c)
		cv.Pixel(cx+y, cy+x, c)
		cv.Pixel(cx-y, cy+x, c)
		cv.Pixel(cx-x, cy+y, c)
		cv.Pixel(cx-x, cy-y, c)
		cv.Pixel(cx-y, cy-x, c)
		cv.Pixel(cx+y, cy-x, c)
		cv.Pixel(cx+x, cy-y, c)
		y++
		if err < 0 {
			err += 2*y + 1
		} else {
			x--
			err += 2*(y-x) + 1
		}
	}
}

// FillRoundRect fills a rectangle with antialiased rounded corners.
func (cv *Canvas) FillRoundRect(x, y, w, h, r int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	r = min(r, min(w, h)/2)
	if r <= 0 {
		cv.FillRect(x, y, w, h, c)
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}

	cv.FillRect(x+r, y, w-2*r, h, c)
	cv.FillRect(x, y+r, r, h-2*r, c)
	cv.FillRect(x+w-r, y+r, r, h-2*r, c)

	x0 := max(x, cv.Clip.X)
	y0 := max(y, cv.Clip.Y)
	x1 := min(x+w, cv.Clip.X+cv.Clip.Width)
	y1 := min(y+h, cv.Clip.Y+cv.Clip.Height)

	for iy := y0; iy < y1; iy++ {
		if iy >= y+r && iy < y+h-r {
			continue
		}
		var dy int
		if iy < y+r {
			dy = iy - (y + r - 1)
		} else {
			dy = iy - (y + h - r)
		}
		for ix := x0; ix < x1; ix++ {
			if ix >= x+r && ix < x+w-r {
				continue
			}
			var dx int
			if ix < x+r {
				dx = ix - (x + r - 1)
			} else {
				dx = ix - (x + w - r)
			}
			alpha := discCoverage(dx, dy, r-1)
			if alpha == 0 {
				continue
			}
			cv.Pixel(ix, iy, Fade(c, alpha*baseAlpha/255))
		}
	}
}

// RoundRect draws a one-pixel outline with antialiased rounded corners.
func (cv *Canvas) RoundRect(x, y, w, h, r int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	r = min(r, min(w, h)/2)
	if r <= 0 {
		cv.Rect(x, y, w, h, c)
		return
	}

	cv.FillRect(x+r, y, w-2*r, 1, c)
	cv.FillRect(x+r, y+h-1, w-2*r, 1, c)
	cv.FillRect(x, y+r, 1, h-2*r, c)
	cv.FillRect(x+w-1, y+r, 1, h-2*r, c)

	corners := [4][2]int{
		{x + r - 1, y + r - 1},
		{x + w - r, y + r - 1},
		{x + r - 1, y + h - r},
		{x + w - r, y + h - r},
	}
	baseAlpha := int(c.A())
	for i, corner := range corners {
		sx, sy := -1, -1
		if i&1 != 0 {
			sx = 1
		}
		if i&2 != 0 {
			sy = 1
		}
		for py := 0; py <= r; py++ {
			for px := 0; px <= r; px++ {
				alpha := discCoverage(px, py, r-1) - discCoverage(px, py, r-2)
				if alpha <= 0 {
					continue
				}
				cv.Pixel(corner[0]+sx*px, corner[1]+sy*py,
					Fade(c, alpha*baseAlpha/255))
			}
		}
	}
}

// FillTriangle fills a triangle. Either winding works.
func (cv *Canvas) FillTriangle(x0, y0, x1, y1, x2, y2 int, c Color) {
	minX := max(min(x0, min(x1, x2)), cv.Clip.X)
	maxX := min(max(x0, max(x1, x2)), cv.Clip.X+cv.Clip.Width-1)
	minY := max(min(y0, min(y1, y2)), cv.Clip.Y)
	maxY := min(max(y0, max(y1, y2)), cv.Clip.Y+cv.Clip.Height-1)

	if area := (x1-x0)*(y2-y0) - (x2-x0)*(y1-y0); area == 0 {
		return
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			w0 := (x1-x0)*(y-y0) - (x-x0)*(y1-y0)
			w1 := (x2-x1)*(y-y1) - (x-x1)*(y2-y1)
			w2 := (x0-x2)*(y-y2) - (x-x2)*(y0-y2)
			if (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0) {
				cv.Pixel(x, y, c)
			}
		}
	}
}

// Triangle draws the three edges of a triangle.
func (cv *Canvas) Triangle(x0, y0, x1, y1, x2, y2 int, c Color) {
	cv.Line(x0, y0, x1, y1, c)
	cv.Line(x1, y1, x2, y2, c)
	cv.Line(x2, y2, x0, y0, c)
}

// Text draws a string and returns the width of the widest line drawn.
func (cv *Canvas) Text(x, y int, text string, c Color) int {
	return cv.TextScaled(x, y, text, c, 1)
}

// TextScaled draws a string larger, and returns the width of the widest line.
func (cv *Canvas) TextScaled(x, y int, text string, c Color, scale int) int {
	scale = max(scale, 1)
	f := cv.face()
	if f != BuiltinFace() {
		f = f.Scaled(scale)
		return f.Draw(cv, x, y+f.Ascent(), text, c)
	}
	return cv.textBitmap(x, y, text, c, scale)
}

// textBitmap is the built-in font.
func (cv *Canvas) textBitmap(x, y int, text string, c Color, scale int) int {
	scale = max(scale, 1)
	penX, penY, widest := x, y, 0

	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, penX-x)
			penX = x
			penY += FontHeight * scale
			continue
		case '\t':
			penX += FontWidth * scale * 4
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
				for col := range FontWidth {
					if bits&(0x80>>col) == 0 {
						continue
					}
					if scale == 1 {
						cv.Pixel(penX+col, penY+row, c)
					} else {
						cv.FillRect(penX+col*scale, penY+row*scale, scale, scale, c)
					}
				}
			}
		}
		penX += FontWidth * scale
	}
	return max(widest, penX-x)
}

// Blit composites another canvas over this one at x, y.
func (cv *Canvas) Blit(x, y int, src *Canvas) {
	if src == nil {
		return
	}
	x0 := max(x, cv.Clip.X)
	y0 := max(y, cv.Clip.Y)
	x1 := min(x+src.Width, cv.Clip.X+cv.Clip.Width)
	y1 := min(y+src.Height, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}

	if cv.Pixels != nil && src.Pixels != nil {
		for sy := y0; sy < y1; sy++ {
			drow := cv.Pixels[sy*cv.Stride : sy*cv.Stride+cv.Width]
			srow := src.Pixels[(sy-y)*src.Stride : (sy-y)*src.Stride+src.Width]
			BlendRow(drow[x0:x1], srow[x0-x:x1-x])
		}
		return
	}
	for by := y0; by < y1; by++ {
		for bx := x0; bx < x1; bx++ {
			cv.Pixel(bx, by, src.At(bx-x, by-y))
		}
	}
}

// TextWidth is the width in pixels that a string takes at scale 1.
func TextWidth(text string) int {
	if f := DefaultFace(); f != BuiltinFace() {
		return f.Width(text)
	}
	return font.TextWidth(text)
}

// TextHeight is how far apart two lines are, in the default face.
func TextHeight() int { return DefaultFace().Height() }

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}