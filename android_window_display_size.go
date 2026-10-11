//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/display"
)

// displaySize is the surface, which on a phone is the display.
func (b *androidWindow) displaySize() (int, int, bool) {
	if b.surface == nil {
		return 0, 0, false
	}
	w, h := b.surface.Width(), b.surface.Height()
	return w, h, w > 0 && h > 0
}
