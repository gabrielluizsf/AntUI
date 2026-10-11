package canvas

// ScaleAlpha multiplies every pixel's alpha, leaving its colour alone. It is
// how a layer painted at full strength is brought down to a colour's own
// transparency before it is composited.
func ScaleAlpha(cv *Canvas, factor float64) {
	if factor >= 1 || factor < 0 || cv.Pixels == nil {
		return
	}
	region := Area{0, 0, cv.Width, cv.Height}
	cv.markDirtyRect(region)
	for i, c := range cv.Pixels {
		cv.Pixels[i] = Fade(c, int(float64(c.A())*factor+0.5))
	}
	cv.markRegion(region)
}
