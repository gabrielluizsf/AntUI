package antui

import "github.com/gabrielluizsf/antui/canvas"

// SetClip narrows drawing to a rectangle, intersected with the window.
func (win *Window) SetClip(a canvas.Area) { win.cv.SetClip(a) }
