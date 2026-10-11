package backend

// Resizable reports whether the window may be resized at all.
func (l Limits) Resizable() bool { return !l.Fixed }
