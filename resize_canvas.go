package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// resizeCanvas grows or shrinks the drawable surface. Backends call it when
// the window system reports a new size.
func (win *Window) ResizeCanvas(width, height int) bool {
	width, height = max(width, 1), max(height, 1)
	if width == win.cv.Width && height == win.cv.Height {
		return true
	}
	cv, err := canvas.NewCanvas(width, height)
	if err != nil {
		return false
	}
	win.cv = cv
	win.width, win.height = width, height

	// The copy of the last frame is the wrong size now, so it cannot be
	// compared against. The buffer behind it is kept: a window that is
	// resized once and then never again would otherwise throw away a
	// megabyte and take it again on the next frame.
	win.shadow = win.shadow[:0]
	win.fullRedraw = true
	win.safe = canvas.Area{}
	return true
}
