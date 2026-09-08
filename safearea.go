package antui

import "github.com/gabrielluizsf/antui/canvas"

// SafeArea is the part of the canvas that nothing else is drawn over.
//
// On a desktop it is the whole canvas and always has been. On a phone it is
// not: an app gets a surface the size of the screen and the system paints
// its status bar, its navigation bar and the notch around the camera **over**
// it — so the app owns every pixel and a strip at each end of them is behind
// something. Anything the user has to be able to read or press belongs
// inside this rectangle; a background belongs outside it, filling the screen.
//
// It is in canvas pixels, and it moves: a rotation changes it, and so does
// the on-screen keyboard coming up.
func (win *Window) SafeArea() canvas.Area {
	full := canvas.Area{X: 0, Y: 0, Width: win.cv.Width, Height: win.cv.Height}
	s := win.safe
	if s.Width <= 0 || s.Height <= 0 {
		return full
	}
	// Clamp rather than trust: the platform reports this in its own
	// coordinates and a stale one after a resize would otherwise send
	// drawing off the end of the canvas.
	if s.X < 0 {
		s.Width += s.X
		s.X = 0
	}
	if s.Y < 0 {
		s.Height += s.Y
		s.Y = 0
	}
	s.Width = min(s.Width, full.Width-s.X)
	s.Height = min(s.Height, full.Height-s.Y)
	if s.Width <= 0 || s.Height <= 0 {
		return full
	}
	return s
}

// setSafeArea is how a backend reports it. A zero rectangle means "the whole
// canvas", which is what every backend that does not know says by saying
// nothing.
func (win *Window) SetSafeArea(a canvas.Area) { win.safe = a }
