package antui

import "github.com/gabrielluizsf/antui/canvas"

// Hovered reports whether the pointer is inside a rectangle.
func (win *Window) Hovered(area canvas.Area) bool {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	return win.mouseX >= x && win.mouseX < x+w &&
		win.mouseY >= y && win.mouseY < y+h
}