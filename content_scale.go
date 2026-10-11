package antui

// ContentScale is how many pixels the window draws for each point the system
// lays it out in: 1 on an ordinary display, 2 on a Mac retina or a Windows
// laptop at 200 per cent, 1.5 at 150.
//
// It is what makes "fixed size" mean the same thing on every machine. A
// system that does not say answers 1 rather than 0, so arithmetic on it is
// always safe.
func (win *Window) ContentScale() float64 {
	if win == nil || win.scale <= 0 {
		return 1
	}
	return win.scale
}
