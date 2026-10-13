package canvas

// MaskRect zeros everything outside a rectangle, keeping only the pixels
// inside it. It is how a background painted into a layer is clipped to part of
// the box — the padding or content box — before the layer is composited.
// It honours the clip on every canvas format.
func (cv *Canvas) MaskRect(a Area) {
	x, y, w, h := a.X, a.Y, a.Width, a.Height
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w || py < y || py >= y+h {
				cv.Put(px, py, 0)
			}
		}
	}
}
