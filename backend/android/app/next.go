//go:build android

package app

// Next is the next event, and false once the activity is destroyed.
//
// It also acknowledges the event before it. Android blocks its UI thread on
// a few of these — WindowDown above all — and the block is released here, on
// the *following* call, or by [App.Release]. That is the contract: an event
// is in hand from the Next that returned it until the next Next, Poll or
// Release, and nothing that the event took away may be touched during it.
func (a *App) Next() (Event, bool) {
	if a.done {
		return None, false
	}
	a.release()
	m, ok := <-a.msgs
	if !ok {
		a.done = true
		return None, false
	}
	return a.take(m), true
}
