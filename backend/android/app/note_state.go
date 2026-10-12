//go:build android

package app

// noteState folds a lifecycle event into the state. It is called from apply,
// on the app's goroutine, before the event is handed to the caller — so an
// app that asks the state while handling Pause is told Paused.
func (a *App) noteState(e Event) {
	switch e {
	case Start:
		a.state = Started
	case Resume:
		a.state = Resumed
	case Pause:
		a.state = Paused
	case Stop:
		a.state = Stopped
	case Destroy:
		a.state = Destroyed
	}
}
