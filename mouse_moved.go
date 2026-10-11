package antui

// MouseMoved reports whether the pointer actually moved during this frame, as
// opposed to sitting still. A widget that lets the mouse move its selection
// (like a calendar's highlighted day) should only claim the selection on
// frames the pointer moved, so the keyboard keeps control otherwise.
func (win *Window) MouseMoved() bool { return win.mouseDX != 0 || win.mouseDY != 0 }
