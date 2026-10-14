package antui

import "github.com/gabrielluizsf/antui/canvas"

// Clicked reports whether a rectangle was clicked or tapped this frame,
// drawing nothing at all.
func (win *Window) Clicked(area canvas.Area) bool {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	return win.widgetClick(widgetID("area", x, y, w, h, ""),
		win.Hovered(area))
}
