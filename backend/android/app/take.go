//go:build android

package app

// take is the half of Next and Poll that is the same: hold the event, fold
// it into the state, and notice the last one.
func (a *App) take(m message) Event {
	a.pending, a.held = m, true
	a.apply(m)
	a.attach(m)
	if m.ev == Destroy {
		a.done = true
		a.release()
		close(a.dead)
	}
	return m.ev
}
