//go:build android

package app

// Config is the configuration the app is running under.
func (a *App) Config() Config {
	if a == nil {
		return Config{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.config
}
