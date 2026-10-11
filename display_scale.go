package antui

// DisplayScale is how many pixels one point of this window covers — the raw
// answer the display gave, and false when it gave none. [Window.ContentScale]
// is the same number with 1 stood in for "did not say", which is what
// arithmetic on a scale wants; the resolution media feature wants the
// difference, because a system that never spoke is one with no resolution to
// report.
//
// The scale is read again while the window runs, so a display the user
// rescales reaches the frame after it; a read that came back with nothing
// leaves the scale the display had said, since the number belongs to the
// screen rather than to the frame that asked.
func (win *Window) DisplayScale() (float64, bool) {
	if win == nil || win.scale <= 0 {
		return 0, false
	}
	return win.scale, true
}
