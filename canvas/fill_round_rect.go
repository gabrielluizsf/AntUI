package canvas

// FillRoundRect fills a rectangle with antialiased corners, each a quarter
// ellipse with the radius [Rounded] gives across and down. Equal radii — what
// almost every rounded rectangle in a UI is — come out as circles, and the
// radii are pulled in to half the box in each direction so a radius wider than
// the box does not turn it inside out.
func (cv *Canvas) FillRoundRect(b Rounded, c Color) {
	x, y, w, h := b.X, b.Y, b.Width, b.Height
	rx, ry := min(b.RX, w/2), min(b.vertical(), h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		cv.FillRect(b.Area, c)
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}

	cv.FillRect(Area{x + rx, y, w - 2*rx, h}, c)
	cv.FillRect(Area{x, y + ry, rx, h - 2*ry}, c)
	cv.FillRect(Area{x + w - rx, y + ry, rx, h - 2*ry}, c)

	// Equal radii are a disc and get the disc's arithmetic, which is the
	// cheaper of the two; only a corner that is genuinely wider one way than
	// the other pays for the ellipse.
	cover := func(dx, dy int) int { return ellipseCoverage(dx, dy, rx-1, ry-1) }
	if rx == ry {
		cover = func(dx, dy int) int { return discCoverage(dx, dy, rx-1) }
	}

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
			alpha := cover(dx, dy)
			if alpha == 0 {
				continue
			}
			cv.Pixel(ix, iy, Fade(c, alpha*baseAlpha/255))
		}
	}
}
