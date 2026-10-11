//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
)

// setFullscreen hides the status bar. The navigation bar stays: taking that
// away is immersive mode, which Android only offers through Java.
func (b *androidWindow) setFullscreen(on bool) bool {
	if on {
		app.SetWindowFlags(app.FlagFullscreen, 0)
	} else {
		app.SetWindowFlags(0, app.FlagFullscreen)
	}
	return true
}
