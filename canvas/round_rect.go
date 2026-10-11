package canvas

// RoundRect draws a one-pixel outline with antialiased rounded corners.
func (cv *Canvas) RoundRect(x, y, w, h, r int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	r = min(r, min(w, h)/2)
	if r <= 0 {
		cv.Rect(x, y, w, h, c)
		return
	}

	cv.FillRect(x+r, y, w-2*r, 1, c)
	cv.FillRect(x+r, y+h-1, w-2*r, 1, c)
	cv.FillRect(x, y+r, 1, h-2*r, c)
	cv.FillRect(x+w-1, y+r, 1, h-2*r, c)

	corners := [4][2]int{
		{x + r - 1, y + r - 1},
		{x + w - r, y + r - 1},
		{x + r - 1, y + h - r},
		{x + w - r, y + h - r},
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
		for py := 0; py <= r; py++ {
			for px := 0; px <= r; px++ {
				alpha := discCoverage(px, py, r-1) - discCoverage(px, py, r-2)
				if alpha <= 0 {
					continue
				}
				cv.Pixel(corner[0]+sx*px, corner[1]+sy*py,
					Fade(c, alpha*baseAlpha/255))
			}
		}
	}
}
