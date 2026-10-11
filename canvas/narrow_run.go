package canvas

// narrowRun writes a whole run in one format, without a call per pixel.
func (cv *Canvas) narrowRun(x0, x1, y int, c Color) {
	switch cv.format {
	case RGB565:
		at := (y*cv.Stride + x0) * 2
		if cv.dither {
			for x := x0; x < x1; x++ {
				packed := to565Dithered(c, x, y)
				cv.narrow[at] = byte(packed)
				cv.narrow[at+1] = byte(packed >> 8)
				at += 2
			}
			return
		}
		packed := To565(c)
		lo, hi := byte(packed), byte(packed>>8)
		for x := x0; x < x1; x++ {
			cv.narrow[at] = lo
			cv.narrow[at+1] = hi
			at += 2
		}
	case Pal8:
		at := y*cv.Stride + x0
		var index byte
		if cv.nearest != nil {
			index = cv.nearest[rgb555(c)]
		} else {
			index = byte(max(cv.PaletteIndex(c), 0))
		}
		row := cv.narrow[at : at+(x1-x0)]
		for i := range row {
			row[i] = index
		}
	}
}
