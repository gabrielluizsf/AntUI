package canvas

// Put writes a colour with no blending, which is what a rasterizer that has
// already decided what it wants needs. It honours the clip.
func (cv *Canvas) Put(x, y int, c Color) {
	if !cv.insideClip(x, y) {
		return
	}
	cv.markDirty(x, y)
	if cv.Pixels != nil {
		// The write happens whatever the baseline says: the canvas has to end
		// up holding what was drawn, because the next draw reads it and a
		// later frame compares against it. The baseline only decides whether
		// the pixel is worth telling the screen about.
		if cv.base != nil && cv.base[y*cv.Width+x] != c {
			cv.markSpan(y, x, x)
		}
		cv.Pixels[y*cv.Stride+x] = c
		return
	}
	if cv.narrow != nil {
		cv.narrowWrite(x, y, c)
	}
}
