package antui

// DroppedFiles is what was dragged onto the window this frame.
func (win *Window) DroppedFiles() []string {
	if win == nil {
		return nil
	}
	return win.dropped
}
