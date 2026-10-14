package antui

import "github.com/gabrielluizsf/antui/canvas"

// ClickedID is [Window.Clicked] for a control that moves: the name decides
// what "the same control" means, rather than where it happens to be.
func (win *Window) ClickedID(name string, x, y, w, h int) bool {
	return win.widgetClick(widgetID("area", 0, 0, 0, 0, name),
		win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h}))
}
