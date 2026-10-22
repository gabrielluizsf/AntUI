package antui

// ClaimModal marks id's popup as the only thing in the window still open to
// the pointer — every other widget's click is refused for as long as the
// claim stands, wherever in the frame's draw order it falls, including one
// drawn earlier in the very same frame and so already past its own hit test
// by the time the popup that will cover it opens. A widget with a popup
// open calls this every frame it stays open; [Window.ReleaseModal] is what
// lets the rest of the window answer again.
func (win *Window) ClaimModal(id uint32) {
	win.modalID = id
}
