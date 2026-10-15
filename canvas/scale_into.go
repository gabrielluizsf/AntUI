package canvas

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
				out.Put(x, y, RGBA(RGBAParams{R: 0, G: 0, B: 0, A: 0}))
				continue
			}
			out.Put(x, y, RGBA(RGBAParams{R: uint8(r / a), G: uint8(g / a), B: uint8(b / a), A: uint8(a / n)}))
		}
	}
}
