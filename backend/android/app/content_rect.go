//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// ContentRect is the part of the window not hidden behind a status bar, a
// navigation bar or a cutout.
func (a *App) ContentRect() ndk.Rect {
	if a == nil {
		return ndk.Rect{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.rect
}
