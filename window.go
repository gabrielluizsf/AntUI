package antui

import (
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// Window is an open window and everything drawn into it. It is not safe for
// concurrent use: like every immediate-mode UI, one goroutine owns the frame
// loop and everything else talks to it through channels.
type Window struct {
	cv     *canvas.Canvas
	native platform

	width, height int
	fullRedraw    bool // forces sending the whole frame
	visible       bool

	windowedWidth, windowedHeight int   // the size to come back out to
	opacity                       uint8 // what was last asked for; 255 = opaque

	limits Limits // what the window may be resized to; see window_size.go

	queue      []Event
	frameStart time.Time

	theme Theme

	windowState
	windowInput
	windowTouch
	windowClock
	windowUI
}
