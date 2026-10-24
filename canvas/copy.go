package canvas

// Copy writes another canvas over this one at x, y as it is, with no
// blending: a pixel of src replaces the one under it, alpha and all. It is
// Blit for a picture known to be opaque — a tile, a background — where the
// blend would come out the same and only costs a look at every pixel. It
// honours the clip.
func (cv *Canvas) Copy(x, y int, src *Canvas) {
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
			end := sy*cv.Stride + cv.Width
			drow := cv.Pixels[sy*cv.Stride : end : end]
			s := (sy-y)*src.Stride - x
			copy(drow[x0:x1], src.Pixels[s+x0:s+x1])
			cv.markRun(sy, x0, x1)
		}
		return
	}
	for by := y0; by < y1; by++ {
		for bx := x0; bx < x1; bx++ {
			cv.Put(bx, by, src.At(bx-x, by-y))
		}
	}
}
