package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// SetSystemFace makes the program draw in the platform's own font, and says
// whether it managed.
//
// It is the whole of what most programs want:
//
//	antui.SetSystemFace(0, win.ContentScale())
//
// and it is deliberately not done automatically. The built-in font is exact
// and identical on every machine, which is what a test card, a screenshot
// comparison and a game with a fixed layout all depend on; changing that
// underneath every program that already exists would break them for a
// nicety.
func SetSystemFace(points, scale float64) bool {
	f, err := UIFace(points, scale)
	if err != nil {
		return false
	}
	canvas.SetDefaultFace(f)
	return true
}
