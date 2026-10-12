//go:build android

package app

// Focused reports whether the window has focus.
func (a *App) Focused() bool {
	if a == nil {
		return false
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.focus
}
