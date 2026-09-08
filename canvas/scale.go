package canvas

// Scaled is a copy of a canvas at a different size.
//
// Every pixel of the result is the average of the block of the original that
// lands on it, weighted by alpha so that the colour of a transparent region
// does not bleed into the edge of what is beside it.
func Scaled(src *Canvas, width, height int) *Canvas {
	if src == nil || width <= 0 || height <= 0 || src.Width <= 0 || src.Height <= 0 {
		return nil
	}
	out, err := NewCanvas(width, height)
	if err != nil {
		return nil
	}
	scaleInto(out, src)
	return out
}

// BlitScaled composites another canvas over this one, fitted into the
// rectangle given rather than at its own size.
func (cv *Canvas) BlitScaled(x, y, width, height int, src *Canvas) {
	if src == nil || width <= 0 || height <= 0 {
		return
	}
	if width == src.Width && height == src.Height {
		cv.Blit(x, y, src)
		return
	}
	if scaled := Scaled(src, width, height); scaled != nil {
		cv.Blit(x, y, scaled)
	}
}

// Fit is the largest width and height with the source's proportions that
// fits inside the box given.
func Fit(srcW, srcH, boxW, boxH int) (width, height int) {
	if srcW <= 0 || srcH <= 0 || boxW <= 0 || boxH <= 0 {
		return 0, 0
	}
	if srcW*boxH < boxW*srcH {
		return srcW * boxH / srcH, boxH
	}
	return boxW, srcH * boxW / srcW
}

// scaleInto is the box filter, writing one canvas into another of any size.
func scaleInto(out, src *Canvas) {
	for y := range out.Height {
		y0 := y * src.Height / out.Height
		y1 := max((y+1)*src.Height/out.Height, y0+1)

		for x := range out.Width {
			x0 := x * src.Width / out.Width
			x1 := max((x+1)*src.Width/out.Width, x0+1)

			var r, g, b, a, n int
			for sy := y0; sy < y1 && sy < src.Height; sy++ {
				for sx := x0; sx < x1 && sx < src.Width; sx++ {
					c := src.At(sx, sy)
					alpha := int(c.A())
					r += int(c.R()) * alpha
					g += int(c.G()) * alpha
					b += int(c.B()) * alpha
					a += alpha
					n++
				}
			}
			if n == 0 || a == 0 {
				out.Put(x, y, RGBA(0, 0, 0, 0))
				continue
			}
			out.Put(x, y, RGBA(uint8(r/a), uint8(g/a), uint8(b/a), uint8(a/n)))
		}
	}
}
