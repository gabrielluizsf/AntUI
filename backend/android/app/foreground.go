//go:build android

package app

// Foreground reports whether the user is looking at the app. It is the one
// question most programs actually want to ask.
func (a *App) Foreground() bool { return a.State() == Resumed }
