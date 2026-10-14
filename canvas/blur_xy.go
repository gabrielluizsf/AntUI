package canvas

// BlurXY is [Canvas.Blur] with a radius for each way. A picture that a
// transform stretched more across than down carries its blur stretched with it,
// and one number cannot say two: radiusX reaches sideways and radiusY up and
// down, and either may be zero for no blur that way.
func (cv *Canvas) BlurXY(area Area, radiusX, radiusY int) {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
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
