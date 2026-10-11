package canvas

// Pixel blends a colour over what is already there, honouring the clip.
func (cv *Canvas) Pixel(x, y int, c Color) {
	if !cv.insideClip(x, y) {
		return
	}
	cv.markDirty(x, y)
	if cv.Pixels == nil {
		if cv.narrow == nil {
			return
		}
		switch c.A() {
		case 0:
			return
		case 255:
			cv.narrowWrite(x, y, c)
		default:
			cv.narrowWrite(x, y, Blend(cv.narrowRead(x, y), c))
		}
		return
	}
	at := y*cv.Stride + x
	// A layer keeps the alpha of what was already there as well as of the
	// colour being drawn, so a shape with a soft edge drawn into one comes out
	// with a soft edge. Everywhere else the destination is a surface that is
	// already solid, and the cheaper Blend that always answers opaque is the
	// same answer.
	var v Color
	if cv.over {
		v = BlendOver(cv.Pixels[at], c)
	} else {
		v = Blend(cv.Pixels[at], c)
	}
	if cv.base != nil && cv.base[y*cv.Width+x] != v {
		cv.markSpan(y, x, x)
	}
	cv.Pixels[at] = v
}
