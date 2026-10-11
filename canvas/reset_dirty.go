package canvas

// ResetDirty clears the dirty bounds, marking the canvas as drawn on nowhere.
// The transformed painting path calls it after wiping its layer, so the
// composite can walk exactly what the widget drew rather than the padded box.
func (cv *Canvas) ResetDirty() { cv.Dirty = Area{} }
