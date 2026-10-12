//go:build android

package app

// drain keeps answering the UI thread after the app's function has returned.
// Without it, a program that simply stops would leave Android blocked in a
// callback and the app would be killed for not responding rather than
// closing.
func (a *App) drain() {
	for a.Alive() {
		if _, ok := a.Next(); !ok {
			return
		}
	}
	a.release()
}
