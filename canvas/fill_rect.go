package canvas

// FillRect fills a rectangle, blending when the colour is translucent.
func (cv *Canvas) FillRect(a Area, c Color) {
	x, y, w, h := a.X, a.Y, a.Width, a.Height
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
	cv.markDirtyRect(Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0})

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
			cv.fillRun(iy, x0, x1, c)
		}
	default:
		for iy := y0; iy < y1; iy++ {
			cv.blendRun(iy, x0, x1, c)
		}
	}
}
