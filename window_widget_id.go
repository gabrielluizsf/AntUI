package antui

// WidgetID is the identity a widget's state is kept under. Two calls with the
// same kind, place and label are the same widget; changing any of them makes
// a new one. The built-in widgets use it, and a template's widgets should
// too, so that a template widget and a built-in one never mistake each other
// for the same thing.
func (win *Window) WidgetID(kind string, x, y, w, h int, label string) uint32 {
	return widgetID(kind, x, y, w, h, label)
}
