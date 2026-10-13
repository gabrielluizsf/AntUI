package canvas

// FillEllipse draws a filled, antialiased ellipse.
func (cv *Canvas) FillEllipse(e Ellipse, c Color) {
	cx, cy, rx, ry := e.X, e.Y, e.RX, e.RY
	if rx <= 0 || ry <= 0 {
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}
	x0 := max(cx-rx-1, cv.Clip.X)
	y0 := max(cy-ry-1, cv.Clip.Y)
	x1 := min(cx+rx+1, cv.Clip.X+cv.Clip.Width-1)
	y1 := min(cy+ry+1, cv.Clip.Y+cv.Clip.Height-1)

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			alpha := ellipseCoverage(x-cx, y-cy, rx, ry)
			if alpha == 0 {
				continue
			}
			cv.Pixel(x, y, Fade(c, alpha*baseAlpha/255))
		}
	}
}
