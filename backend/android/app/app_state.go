//go:build android

package app

// State is where the activity is now.
func (a *App) State() State {
	if a == nil {
		return Destroyed
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}
