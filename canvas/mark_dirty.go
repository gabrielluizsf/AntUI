package canvas

// markDirty widens the dirty bounds to take the point in.
func (cv *Canvas) markDirty(x, y int) {
	if cv.Dirty.Width <= 0 {
		cv.Dirty = Area{X: x, Y: y, Width: 1, Height: 1}
		return
	}
	x1, y1 := cv.Dirty.X+cv.Dirty.Width, cv.Dirty.Y+cv.Dirty.Height
	if x < cv.Dirty.X {
		cv.Dirty.X = x
	}
	if x1 <= x {
		x1 = x + 1
	}
	if y < cv.Dirty.Y {
		cv.Dirty.Y = y
	}
	if y1 <= y {
		y1 = y + 1
	}
	cv.Dirty.Width, cv.Dirty.Height = x1-cv.Dirty.X, y1-cv.Dirty.Y
}
