//go:build android

package app

// Density is the screen's density in dots per inch. It is [Config.Density],
// kept as its own call because it is what the window layer asks for on every
// configuration change.
func (a *App) Density() int {
	if a == nil {
		return 0
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.config.Density
}
