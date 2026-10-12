//go:build android

package jni

// Activity is the app's android.app.Activity, as a global reference. It is
// the starting point for almost everything: a Context, a window, a place to
// ask for a system service.
func Activity() Object {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return Object{ref: activity}
}
