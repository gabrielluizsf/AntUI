package template

import (
	"sync"

	"github.com/gabrielluizsf/antui/canvas"
)

// The registered icons, and the pictures already painted from them.
//
// The lock is because a window paints on its own goroutine while the program
// that owns it may replace an icon at any moment; the map is read on every
// widget of every frame, so the read side is the one that matters and takes the
// lock only for as long as it takes to read one pointer.
var (
	iconMu sync.RWMutex
	icons  = map[string]*iconEntry{}

	// Painted icons, keyed by what they were painted at. A drawing remembers one
	// size and colour of itself, which is not enough for an interface that holds
	// a list of buttons: two of them in the same state share a picture, but the
	// same button hovered and not hovered alternates between two colours every
	// frame, and a drawing that could only remember one of them would redraw the
	// icon on every frame forever. So the pictures are kept here instead, where a
	// key covers all three.
	iconPaintMu sync.RWMutex
	iconPaint   = map[iconPaintKey]*canvas.Canvas{}
)
