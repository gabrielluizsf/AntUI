//go:build android

package app

// Poll is Next for a caller that cannot afford to wait: it returns the next
// event if one is already there, and None with false if none is. The
// acknowledgement rule is the same — the event before is released here, so a
// frame loop that polls once a frame holds the UI thread for at most a frame.
//
// A false result means "nothing right now", not "the app is over".
// [App.Alive] is what says that.
func (a *App) Poll() (Event, bool) {
	if a.done {
		return None, false
	}
	a.release()
	select {
	case m, ok := <-a.msgs:
		if !ok {
			a.done = true
			return None, false
		}
		return a.take(m), true
	default:
		return None, false
	}
}
