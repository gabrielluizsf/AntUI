package antui

import "github.com/gabrielluizsf/antui/canvas"

// ClickedID is [Window.Clicked] for a control that moves: the name decides
// what "the same control" means, rather than where it happens to be.
func (win *Window) ClickedID(name string, area canvas.Area) bool {
	return win.widgetClick(widgetID("area", 0, 0, 0, 0, name),
		win.Hovered(area))
}
