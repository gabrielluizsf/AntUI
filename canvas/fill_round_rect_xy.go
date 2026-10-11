package canvas

// FillRoundRectXY fills a rectangle whose corners are quarter ellipses with
// horizontal radius rx and vertical radius ry.
func (cv *Canvas) FillRoundRectXY(x, y, w, h, rx, ry int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		cv.FillRect(x, y, w, h, c)
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}

	cv.FillRect(x+rx, y, w-2*rx, h, c)
	cv.FillRect(x, y+ry, rx, h-2*ry, c)
	cv.FillRect(x+w-rx, y+ry, rx, h-2*ry, c)

	x0 := max(x, cv.Clip.X)
	y0 := max(y, cv.Clip.Y)
	x1 := min(x+w, cv.Clip.X+cv.Clip.Width)
	y1 := min(y+h, cv.Clip.Y+cv.Clip.Height)

	for iy := y0; iy < y1; iy++ {
		if iy >= y+ry && iy < y+h-ry {
			continue
		}
		var dy int
		if iy < y+ry {
			dy = iy - (y + ry - 1)
		} else {
			dy = iy - (y + h - ry)
		}
		for ix := x0; ix < x1; ix++ {
			if ix >= x+rx && ix < x+w-rx {
				continue
			}
			var dx int
			if ix < x+rx {
				dx = ix - (x + rx - 1)
			} else {
				dx = ix - (x + w - rx)
			}
			alpha := ellipseCoverage(dx, dy, rx-1, ry-1)
			if alpha == 0 {
				continue
			}
			cv.Pixel(ix, iy, Fade(c, alpha*baseAlpha/255))
		}
	}
}
