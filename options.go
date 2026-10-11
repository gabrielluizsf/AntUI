package antui

import (
	"github.com/gabrielluizsf/antui/backend"
)

// Options is everything a window can be opened with. The zero value plus a
// size is exactly what Open gives: a resizable window with no bounds.
type Options = backend.Options
