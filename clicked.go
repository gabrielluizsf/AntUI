package antui

// Clicked reports whether a rectangle was clicked or tapped this frame,
// drawing nothing at all.
func (win *Window) Clicked(x, y, w, h int) bool {
	return win.widgetClick(widgetID("area", x, y, w, h, ""),
		win.Hovered(x, y, w, h))
}
