package canvas

type fillRoundRectCorners struct {
	cv            *Canvas
	c             Color
	x, y, w, h    int
	rx, ry, alpha int
}

func (s fillRoundRectCorners) draw() {
	// Equal radii use the cheaper disc arithmetic; unequal corners need the ellipse.
	cover := func(dx, dy int) int { return ellipseCoverage(dx, dy, s.rx-1, s.ry-1) }
	if s.rx == s.ry {
		cover = func(dx, dy int) int { return discCoverage(dx, dy, s.rx-1) }
	}

	x0 := max(s.x, s.cv.Clip.X)
	y0 := max(s.y, s.cv.Clip.Y)
	x1 := min(s.x+s.w, s.cv.Clip.X+s.cv.Clip.Width)
	y1 := min(s.y+s.h, s.cv.Clip.Y+s.cv.Clip.Height)
	for iy := y0; iy < y1; iy++ {
		if iy >= s.y+s.ry && iy < s.y+s.h-s.ry {
			continue
		}
		dy := iy - (s.y + s.ry - 1)
		if iy >= s.y+s.ry {
			dy = iy - (s.y + s.h - s.ry)
		}
		for ix := x0; ix < x1; ix++ {
			if ix >= s.x+s.rx && ix < s.x+s.w-s.rx {
				continue
			}
			dx := ix - (s.x + s.rx - 1)
			if ix >= s.x+s.rx {
				dx = ix - (s.x + s.w - s.rx)
			}
			alpha := cover(dx, dy)
			if alpha > 0 {
				s.cv.Pixel(ix, iy, Fade(s.c, alpha*s.alpha/255))
			}
		}
	}
}
