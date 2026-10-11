package antui

import (
	"github.com/gabrielluizsf/antui/backend"
)

// Key identifies a physical key, independent of the text it produces.
// Printable keys carry their own uppercase ASCII code, so KeyA is 'A' and
// Key1 is '1'; everything else is numbered from 256 up.
type Key = backend.Key
