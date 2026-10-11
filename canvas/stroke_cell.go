package canvas

func (cv *Canvas) strokeCell(x, y, scale int, c Color, bold bool) {
	if scale <= 1 {
		cv.Pixel(x, y, c)
		if bold {
			cv.Pixel(x+1, y, c)
		}
		return
	}
	cv.FillRect(x, y, scale, scale, c)
	if bold {
		cv.FillRect(x+scale, y, scale, scale, c)
	}
}
