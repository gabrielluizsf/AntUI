package canvas

// Changes walks the rectangles that really moved in the last committed frame.
// The window uploads their bounding box; a backend that damages per box —
// Wayland, Android — walks this instead and sends only the pieces that moved.
func (cv *Canvas) Changes(yield func(area Area) bool) {
	if cv.change.Width <= 0 {
		return
	}
	x, y := cv.change.X, cv.change.Y
	w, h := cv.change.Width, cv.change.Height
	for dy := range h {
		if !yield(Area{X: x, Y: y + dy, Width: w, Height: 1}) {
			return
		}
	}
}
