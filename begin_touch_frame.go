package antui

// beginTouchFrame is the per-frame reset: the fingers that lifted last frame
// are gone, and nothing has moved yet.
func (win *Window) beginTouchFrame() {
	kept := win.touches[:0]
	for _, t := range win.touches {
		if t.Ended {
			continue
		}
		t.Began = false
		t.DX, t.DY = 0, 0
		kept = append(kept, t)
	}
	win.touches = kept
}
