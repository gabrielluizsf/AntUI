package antui

import "github.com/gabrielluizsf/antui/canvas"

// Copy writes another canvas over the window with no blending. See
// canvas.Canvas.Copy: it is Blit for a picture known to be opaque.
func (win *Window) Copy(x, y int, src *canvas.Canvas) { win.cv.Copy(x, y, src) }
