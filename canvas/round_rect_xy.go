package canvas

// RoundRectXY draws a one-pixel outline with antialiased quarter-ellipse
// corners with horizontal radius rx and vertical radius ry.
func (cv *Canvas) RoundRectXY(x, y, w, h, rx, ry int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		cv.Rect(x, y, w, h, c)
		return
	}

	cv.FillRect(x+rx, y, w-2*rx, 1, c)
	cv.FillRect(x+rx, y+h-1, w-2*rx, 1, c)
	cv.FillRect(x, y+ry, 1, h-2*ry, c)
	cv.FillRect(x+w-1, y+ry, 1, h-2*ry, c)

	corners := [4][2]int{
		{x + rx - 1, y + ry - 1},
		{x + w - rx, y + ry - 1},
		{x + rx - 1, y + h - ry},
		{x + w - rx, y + h - ry},
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
		for py := 0; py <= ry; py++ {
			for px := 0; px <= rx; px++ {
				alpha := ellipseCoverage(px, py, rx-1, ry-1) - ellipseCoverage(px, py, rx-2, ry-2)
				if alpha <= 0 {
					continue
				}
				cv.Pixel(corner[0]+sx*px, corner[1]+sy*py,
					Fade(c, alpha*baseAlpha/255))
			}
		}
	}
}
