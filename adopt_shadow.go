package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// adoptShadow takes the frame that is about to be sent whole as the one to
// compare against from now on, and gives the canvas the copy to check writes
// with. The pixels of the canvas are what the screen will hold, so the next
// frame starts from them.
func (win *Window) adoptShadow(cv *canvas.Canvas) {
	count := cv.Width * cv.Height
	if cap(win.shadow) < count {
		win.shadow = make([]canvas.Color, count)
	} else {
		win.shadow = win.shadow[:count]
	}
	copy(win.shadow, cv.Pixels)
	cv.Compare(win.shadow)
}
