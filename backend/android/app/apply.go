//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// apply folds an event into the app's state before the caller sees it — so
// that on WindowDown the window is already gone, and on WindowUp it is
// already there.
func (a *App) apply(m message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.noteState(m.ev)
	switch m.ev {
	case WindowUp:
		a.window = ndk.Wrap(m.ptr)
	case WindowDown:
		a.window = nil
	case FocusGained:
		a.focus = true
	case FocusLost:
		a.focus = false
	}
}
