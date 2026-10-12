//go:build android

package app

func current() *App {
	stateMu.Lock()
	defer stateMu.Unlock()
	return theApp
}
