package antui

// SetClip narrows drawing to a rectangle.
func (win *Window) SetClip(x, y, w, h int) { win.cv.SetClip(x, y, w, h) }
