package canvas

// SetDither turns ordered dithering on or off for a narrow canvas. It is on
// by default, which is where it costs nothing and is any use.
func (cv *Canvas) SetDither(on bool) { cv.dither = on }
