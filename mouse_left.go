package antui

import (
	"github.com/gabrielluizsf/antui/backend"
)

// The mouse buttons.
const (
	MouseLeft   = backend.MouseLeft
	MouseRight  = backend.MouseRight
	MouseMiddle = backend.MouseMiddle
	// mouseCount bounds the mouse-state tables. It comes after the buttons.
	mouseCount = 3
)
