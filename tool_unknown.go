package antui

import (
	"github.com/gabrielluizsf/antui/backend"
)

// The tools. A device that does not say reports ToolUnknown rather than
// guessing at a finger.
const (
	ToolUnknown = backend.ToolUnknown
	ToolFinger  = backend.ToolFinger
	ToolStylus  = backend.ToolStylus
	ToolMouse   = backend.ToolMouse
	ToolEraser  = backend.ToolEraser // the far end of a stylus, used to rub out
)
