//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/display"
)

// displayRefresh is the panel's rate, asked once and remembered: it comes
// from Java, and a round trip to the machine for every caller of
// Window.DisplayRefresh would be a poor way to answer a number that changes
// about as often as the hardware does.
//
// It is rounded, because this is where the honest float has to become the
// integer antui's API promised. 59.94 is 60.
func (b *androidWindow) displayRefresh() int {
	if b.refresh == 0 {
		b.refresh = -1 // asked and unanswerable, so do not ask again
		if hz, err := display.Refresh(); err == nil && hz > 0 {
			b.refresh = int(hz + 0.5)
		}
	}
	if b.refresh < 0 {
		return 0
	}
	return b.refresh
}
