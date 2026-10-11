package antui

import (
	"github.com/gabrielluizsf/antui/backend"
)

// Limits are what a window may be resized to. The zero value is a window
// with no bounds at all, which is what a window is without them — so
// Limits{MinWidth: 640, MinHeight: 480} says only what it looks like it says
// and leaves everything else alone.
type Limits = backend.Limits
