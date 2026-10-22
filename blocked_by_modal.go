package antui

// blockedByModal reports whether id is refused a click because some other
// widget's popup has claimed the screen. Every widget but the claim's own is
// refused outright, inside the popup's own rectangle included — a click
// that lands there is answered by the popup's own interaction code, which
// reads the pointer directly rather than going through [Window.WidgetClick]
// at all, not by whatever this window happened to draw under it.
func (win *Window) blockedByModal(id uint32) bool {
	return win.modalID != 0 && win.modalID != id
}
