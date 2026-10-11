package antui

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// windowTouch is the text buffer and the touch and gesture state.
type windowTouch struct {
	text strings.Builder
	// dropped is what was dragged onto the window this frame.
	dropped []string
	// touchFirst is a platform with nothing but a touchscreen, and sawTouch
	// is a finger having actually arrived. See Window.TouchScreen.
	touchFirst bool
	sawTouch   bool
	// touches is every finger on the screen; see touch.go.
	touches []Touch
	// gesture is what those fingers are in the middle of doing; see
	// gesture.go.
	gesture gestureState
	// safe is the part of the canvas the system is not covering; see
	// safearea.go. Zero means all of it.
	safe canvas.Area
}
