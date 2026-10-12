//go:build android

package app

// Alive reports whether the activity still exists. It goes false once
// Destroy has been taken, and stays false.
func (a *App) Alive() bool { return a != nil && !a.done }
