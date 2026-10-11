package antui

// Hovered reports whether the pointer is inside a rectangle.
func (win *Window) Hovered(x, y, w, h int) bool {
	return win.mouseX >= x && win.mouseX < x+w &&
		win.mouseY >= y && win.mouseY < y+h
}
