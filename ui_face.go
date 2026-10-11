package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// UIFace is [SystemFace] at the size the platform draws its own interface
// at, scaled for the display.
//
// Points, not pixels: the size an interface is specified in is a physical
// one, and the number of pixels it comes to depends on the screen. A program
// with a window passes win.ContentScale(); zero points means the platform's
// own default, which differs by platform and is why it is not a constant
// here.
func UIFace(points, scale float64) (*canvas.Face, error) {
	if points <= 0 {
		points = defaultUIPoints
	}
	if scale <= 0 {
		scale = 1
	}
	return SystemFace(points * scale)
}
