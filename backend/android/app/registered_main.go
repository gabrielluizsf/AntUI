//go:build android

package app

func registeredMain() func(*App) {
	mainMu.Lock()
	defer mainMu.Unlock()
	if mainFn != nil {
		return mainFn
	}
	return fallback
}
