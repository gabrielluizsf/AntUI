package canvas

// Blit composites another canvas over this one at x, y.
func (cv *Canvas) Blit(x, y int, src *Canvas) {
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
	cv.markDirtyRect(Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0})

	if cv.Pixels != nil && src.Pixels != nil {
		for sy := y0; sy < y1; sy++ {
			drow := cv.Pixels[sy*cv.Stride : sy*cv.Stride+cv.Width]
			srow := src.Pixels[(sy-y)*src.Stride : (sy-y)*src.Stride+src.Width]
			BlendRow(drow[x0:x1], srow[x0-x:x1-x])
			cv.markRun(sy, x0, x1)
		}
		return
	}
	for by := y0; by < y1; by++ {
		for bx := x0; bx < x1; bx++ {
			cv.Pixel(bx, by, src.At(bx-x, by-y))
		}
	}
}
