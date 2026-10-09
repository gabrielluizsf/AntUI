package canvas

// Blur blurs the rectangle in place. Three box passes stand in for a Gaussian,
// which is what CSS blur radii are measured in; the colour is premultiplied
// while it moves so a transparent pixel cannot drag its black into the edges.
// It works on ARGB32 only.
func (cv *Canvas) Blur(x, y, w, h, radius int) {
	cv.BlurXY(x, y, w, h, radius, radius)
}

// BlurXY is [Canvas.Blur] with a radius for each way. A picture that a
// transform stretched more across than down carries its blur stretched with it,
// and one number cannot say two: radiusX reaches sideways and radiusY up and
// down, and either may be zero for no blur that way.
func (cv *Canvas) BlurXY(x, y, w, h, radiusX, radiusY int) {
	if cv.Pixels == nil || (radiusX <= 0 && radiusY <= 0) || w <= 0 || h <= 0 {
		return
	}
	x0, y0, x1, y1 := cv.filterBounds(x, y, w, h)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	region := Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0}
	cv.markDirtyRect(region)
	rw, rh := x1-x0, y1-y0

	// The picture is moved between two float buffers rather than a new one per
	// pass. Six passes over a card's worth of pixels is a lot of memory to ask
	// the collector for, and the arithmetic is the same either way.
	front, back := cv.blurBuffers(rw * rh * 4)
	for row := range rh {
		for col := range rw {
			c := cv.Pixels[(y0+row)*cv.Stride+x0+col]
			a := float64(c.A()) / 255
			i := (row*rw + col) * 4
			front[i] = float64(c.R()) * a
			front[i+1] = float64(c.G()) * a
			front[i+2] = float64(c.B()) * a
			front[i+3] = a
		}
	}
	for range 3 {
		back = boxBlur(front, back, rw, rh, radiusX, true)
		front, back = back, front
		back = boxBlur(front, back, rw, rh, radiusY, false)
		front, back = back, front
	}
	buf := front
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
	cv.markRegion(region)
}

// blurBuffers is the pair of float buffers a blur moves its picture between,
// sized to n floats and kept on the canvas so the next blur of the same
// canvas — next frame, most likely — finds them already there.
func (cv *Canvas) blurBuffers(n int) (front, back []float64) {
	if cap(cv.blurFront) < n {
		cv.blurFront = make([]float64, n)
		cv.blurBack = make([]float64, n)
	}
	return cv.blurFront[:n], cv.blurBack[:n]
}

// boxBlur runs one moving-average pass, horizontally or vertically, clamping
// at the edges so the picture does not darken as it approaches them. The result
// goes into dst, which the caller owns and hands back on the next pass.
func boxBlur(src, dst []float64, w, h, r int, horizontal bool) []float64 {
	span := float64(2*r + 1)
	if horizontal {
		for row := range h {
			boxBlurLine(src, dst, row*w*4, 4, w, r, span)
		}
		return dst
	}
	for col := range w {
		boxBlurLine(src, dst, col*4, w*4, h, r, span)
	}
	return dst
}

// boxBlurLine is one row or one column of a pass. The four channels are
// averaged side by side rather than one after another — each keeps its own
// running sum and its own order of additions, so the numbers are the same ones
// four separate walks would have produced, and the line is read once instead of
// four times.
//
// The window reaches past both ends of the line and has to be clamped there,
// which is two samples of bookkeeping for every pixel of a shadow. The middle
// of the line, where the window stays inside, skips it.
func boxBlurLine(src, dst []float64, at, stride, count, r int, span float64) {
	var sr, sg, sb, sa float64
	for k := -r; k <= r; k++ {
		i := at + clampInt(k, 0, count-1)*stride
		sr += src[i]
		sg += src[i+1]
		sb += src[i+2]
		sa += src[i+3]
	}

	// ends is the run of samples whose window reaches past the start of the
	// line, mid is where it stops reaching past the end, and everything between
	// them is a sample both of whose ends are inside the line: n-r is a real
	// sample from n=r on, and n+r+1 is one until n = count-r-2.
	ends := min(r, count)
	mid := min(max(count-r-1, ends), count)

	for n := range ends {
		o := at + n*stride
		dst[o] = sr / span
		dst[o+1] = sg / span
		dst[o+2] = sb / span
		dst[o+3] = sa / span
		lo, hi := at+clampInt(n-r, 0, count-1)*stride, at+clampInt(n+r+1, 0, count-1)*stride
		sr = sr - src[lo] + src[hi]
		sg = sg - src[lo+1] + src[hi+1]
		sb = sb - src[lo+2] + src[hi+2]
		sa = sa - src[lo+3] + src[hi+3]
	}
	for n := ends; n < mid; n++ {
		o := at + n*stride
		dst[o] = sr / span
		dst[o+1] = sg / span
		dst[o+2] = sb / span
		dst[o+3] = sa / span
		lo, hi := o-r*stride, o+(r+1)*stride
		sr = sr - src[lo] + src[hi]
		sg = sg - src[lo+1] + src[hi+1]
		sb = sb - src[lo+2] + src[hi+2]
		sa = sa - src[lo+3] + src[hi+3]
	}
	for n := mid; n < count; n++ {
		o := at + n*stride
		dst[o] = sr / span
		dst[o+1] = sg / span
		dst[o+2] = sb / span
		dst[o+3] = sa / span
		lo, hi := at+clampInt(n-r, 0, count-1)*stride, at+clampInt(n+r+1, 0, count-1)*stride
		sr = sr - src[lo] + src[hi]
		sg = sg - src[lo+1] + src[hi+1]
		sb = sb - src[lo+2] + src[hi+2]
		sa = sa - src[lo+3] + src[hi+3]
	}
}
