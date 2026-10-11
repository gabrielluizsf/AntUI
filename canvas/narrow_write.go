package canvas

func (cv *Canvas) narrowWrite(x, y int, c Color) {
	at := y*cv.Stride + x
	switch {
	case cv.format == RGB565:
		var packed uint16
		if cv.dither {
			packed = to565Dithered(c, x, y)
		} else {
			packed = To565(c)
		}
		cv.narrow[at*2] = byte(packed)
		cv.narrow[at*2+1] = byte(packed >> 8)
	case cv.format == Pal8 && cv.nearest != nil:
		cv.narrow[at] = cv.nearest[rgb555(c)]
	case cv.format == Pal8:
		cv.narrow[at] = byte(max(cv.PaletteIndex(c), 0))
	}
}
