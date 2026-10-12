//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// OnNewIntent registers what to do when an intent reaches an activity that
// is already running. The Intent handed to f is a **global reference f
// owns** and must give back.
//
// It is called on a goroutine of its own, not on the UI thread and not on
// the app's.
func OnNewIntent(f func(jni.Object)) {
	intentMu.Lock()
	intentFn = f
	intentMu.Unlock()
}
