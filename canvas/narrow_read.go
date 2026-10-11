package canvas

// narrowRead and narrowWrite are what every narrow path is built out of.
// Neither checks the clip: the callers do.
func (cv *Canvas) narrowRead(x, y int) Color {
	at := y*cv.Stride + x
	switch {
	case cv.format == RGB565:
		lo := uint16(cv.narrow[at*2])
		hi := uint16(cv.narrow[at*2+1])
		return From565(lo | hi<<8)
	case cv.format == Pal8 && cv.palette != nil:
		return cv.palette[cv.narrow[at]]
	}
	return 0
}
