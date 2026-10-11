package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// dirtyRegion is the rectangle to send: the pixels of this frame that really
// changed, which the canvas knows as each one is written because the canvas is
// comparing against the frame already on the screen. A frame that redraws the
// same picture changes nothing and sends nothing — there is no scan of the
// canvas left to find that out, because nothing has to look.
func (win *Window) dirtyRegion() canvas.Area {
	cv := win.cv
	whole := canvas.Area{X: 0, Y: 0, Width: cv.Width, Height: cv.Height}
	if win.fullRedraw || len(win.shadow) != cv.Width*cv.Height {
		win.fullRedraw = false
		win.adoptShadow(cv)
		return whole
	}
	a := cv.PresentChanges()
	return a
}
