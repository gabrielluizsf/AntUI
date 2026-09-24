package canvas

// Blur blurs the rectangle in place. Three box passes stand in for a Gaussian,
// which is what CSS blur radii are measured in; the colour is premultiplied
// while it moves so a transparent pixel cannot drag its black into the edges.
// It works on ARGB32 only.
func (cv *Canvas) Blur(x, y, w, h, radius int) {
	if cv.Pixels == nil || radius <= 0 || w <= 0 || h <= 0 {
		return
	}
	x0, y0, x1, y1 := cv.filterBounds(x, y, w, h)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	cv.markDirtyRect(Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0})
	rw, rh := x1-x0, y1-y0
	buf := make([]float64, rw*rh*4)
	for row := range rh {
		for col := range rw {
			c := cv.Pixels[(y0+row)*cv.Stride+x0+col]
			a := float64(c.A()) / 255
			i := (row*rw + col) * 4
			buf[i] = float64(c.R()) * a
			buf[i+1] = float64(c.G()) * a
			buf[i+2] = float64(c.B()) * a
			buf[i+3] = a
		}
	}
	for range 3 {
		buf = boxBlur(buf, rw, rh, radius, true)
		buf = boxBlur(buf, rw, rh, radius, false)
	}
	for row := range rh {
		for col := range rw {
			i := (row*rw + col) * 4
			a := buf[i+3]
			var r, g, b float64
			if a > 0.0001 {
				r, g, b = buf[i]/a, buf[i+1]/a, buf[i+2]/a
			}
			cv.Pixels[(y0+row)*cv.Stride+x0+col] = RGBA(clampByte(r), clampByte(g), clampByte(b), clampByte(a*255))
		}
	}
}

// boxBlur runs one moving-average pass, horizontally or vertically, clamping
// at the edges so the picture does not darken as it approaches them.
func boxBlur(src []float64, w, h, r int, horizontal bool) []float64 {
	dst := make([]float64, len(src))
	span := float64(2*r + 1)
	if horizontal {
		for row := range h {
			base := row * w * 4
			for ch := range 4 {
				var sum float64
				for k := -r; k <= r; k++ {
					sum += src[base+clampInt(k, 0, w-1)*4+ch]
				}
				for col := range w {
					dst[base+col*4+ch] = sum / span
					sum -= src[base+clampInt(col-r, 0, w-1)*4+ch]
					sum += src[base+clampInt(col+r+1, 0, w-1)*4+ch]
				}
			}
		}
		return dst
	}
	for col := range w {
		for ch := range 4 {
			var sum float64
			for k := -r; k <= r; k++ {
				sum += src[clampInt(k, 0, h-1)*w*4+col*4+ch]
			}
			for row := range h {
				dst[row*w*4+col*4+ch] = sum / span
				sum -= src[clampInt(row-r, 0, h-1)*w*4+col*4+ch]
				sum += src[clampInt(row+r+1, 0, h-1)*w*4+col*4+ch]
			}
		}
	}
	return dst
}
