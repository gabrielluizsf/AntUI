package canvas

type roundRectCorners struct {
	cv      *Canvas
	c       Color
	rx, ry  int
	corners [4][2]int
}

func (s roundRectCorners) draw() {
	// Equal radii use the cheaper disc arithmetic; unequal corners need the ellipse.
	cover := func(px, py, shrink int) int {
		return ellipseCoverage(px, py, s.rx-shrink, s.ry-shrink)
	}
	if s.rx == s.ry {
		cover = func(px, py, shrink int) int {
			return discCoverage(px, py, s.rx-shrink)
		}
	}

	baseAlpha := int(s.c.A())
	for i, corner := range s.corners {
		sx, sy := -1, -1
		if i&1 != 0 {
			sx = 1
		}
		if i&2 != 0 {
			sy = 1
		}
		for py := 0; py <= s.ry; py++ {
			for px := 0; px <= s.rx; px++ {
				alpha := cover(px, py, 1) - cover(px, py, 2)
				if alpha <= 0 {
					continue
				}
				s.cv.Pixel(corner[0]+sx*px, corner[1]+sy*py,
					Fade(s.c, alpha*baseAlpha/255))
			}
		}
	}
}
