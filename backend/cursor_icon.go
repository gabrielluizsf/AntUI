package backend

import "github.com/gabrielluizsf/antui/canvas"

// CursorIcon is a custom mouse pointer: the picture, and which pixel of it
// is the point a click actually lands on — a crosshair's centre, an
// arrow's tip.
type CursorIcon struct {
	Image      *canvas.Canvas
	HotX, HotY int
}
