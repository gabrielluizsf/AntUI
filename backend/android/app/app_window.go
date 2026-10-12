//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Window is the surface to draw on, or nil when there is none. It is only
// valid between WindowUp and WindowDown, and the whole point of Next's
// contract is that "between" means what it says.
func (a *App) Window() *ndk.Window {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.window
}
