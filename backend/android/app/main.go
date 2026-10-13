//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// Main registers the function that is the app. Call it from an init
// function: a shared library's main is never run, so init is the only place
// that happens before Android calls in.
//
//	func init() { app.Main(run) }
//
// Calling it twice replaces the first, which is only ever a mistake, so the
// second call says so in the log.
func Main(f func(*App)) {
	mainMu.Lock()
	defer mainMu.Unlock()
	if mainFn != nil {
		ndk.Warnf("app.Main called twice; the second function wins")
	}
	mainFn = f
}
