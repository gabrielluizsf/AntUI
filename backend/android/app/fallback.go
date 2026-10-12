//go:build android

package app

// Fallback registers f as the app only if nothing registers one with [Main].
//
// It is what the antuiapk command injects, pointing at the program's own
// main, so that a program written for a desktop runs on a phone with no
// Android in its source at all. Which init runs first does not matter: the
// choice is made when Android calls in, not when the registration happens.
func Fallback(f func()) {
	mainMu.Lock()
	defer mainMu.Unlock()
	fallback = func(*App) { f() }
}
