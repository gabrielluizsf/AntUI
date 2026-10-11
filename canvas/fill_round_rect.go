package canvas

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
