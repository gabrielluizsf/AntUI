package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// windowState is the window's durable state: what was last presented, asked
// for, or overridden.
type windowState struct {
	// shadow is the frame that was last presented, which the canvas checks
	// every write against: what the screen already holds is what does not
	// have to be sent again. It is grown once and kept between resizes.
	shadow []canvas.Color

	shouldClose bool
	presented   bool // whether the last End reached the display
	fullscreen  bool // what was last asked for

	scale       float64 // pixels per point, 0 until the backend says
	scaleForced bool    // a SetDisplayScale override; see window_system.go

	// dark is the color scheme the system paints its own interface in,
	// darkKnown says whether the platform had an answer to give, and
	// darkForced is a SetSystemDark override standing in for both. See
	// window_system.go.
	dark, darkKnown, darkForced bool

	droppedEvents int

	// frontmost are draws that run at the very end of the frame, after every
	// widget painted in flow order, so an open dropdown or calendar stays on
	// top of whatever the layout drew under it.
	frontmost []func()
}
