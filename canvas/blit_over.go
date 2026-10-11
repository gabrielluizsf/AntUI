package canvas

// BlitOver composites src over this canvas honouring both alphas, so a soft
// shadow can land on a recording layer that is still transparent. Blit blends
// the same pixels but answers fully opaque, which makes it the wrong tool
// where the destination has not been painted yet. Over an opaque canvas the
// two produce the same colours.
func (cv *Canvas) BlitOver(x, y int, src *Canvas) {
	if src == nil {
		return
	}
	x0 := max(x, cv.Clip.X)
	y0 := max(y, cv.Clip.Y)
	x1 := min(x+src.Width, cv.Clip.X+cv.Clip.Width)
	y1 := min(y+src.Height, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	if cv.Pixels != nil && src.Pixels != nil {
		cv.markDirtyRect(Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0})
		for sy := y0; sy < y1; sy++ {
			drow := cv.Pixels[sy*cv.Stride : sy*cv.Stride+cv.Width]
			srow := src.Pixels[(sy-y)*src.Stride : (sy-y)*src.Stride+src.Width]
			for i := x0; i < x1; i++ {
				if c := srow[i-x]; c != 0 {
					drow[i] = BlendOver(drow[i], c)
				}
			}
			cv.markRun(sy, x0, x1)
		}
		return
	}
	for by := y0; by < y1; by++ {
		for bx := x0; bx < x1; bx++ {
			if c := src.At(bx-x, by-y); c != 0 {
				cv.Pixel(bx, by, c)
			}
		}
	}
}
