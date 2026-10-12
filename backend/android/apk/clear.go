package apk

import (
	_ "image/jpeg"

	"github.com/gabrielluizsf/antui/canvas"
)

// clear empties a canvas to nothing at all.
//
// Not Canvas.Clear, which forces the alpha to 255 — deliberately, because it
// is for clearing a *window* and a frame on screen is opaque. An icon is the
// opposite: what is not drawn has to be see-through, and a rounded corner
// filled with opaque black is a black square.
func clear(cv *canvas.Canvas) {
	for y := range cv.Height {
		for x := range cv.Width {
			cv.Put(x, y, canvas.Transparent)
		}
	}
}
