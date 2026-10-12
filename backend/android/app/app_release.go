//go:build android

package app

// release lets the UI thread go, if it was waiting on the event in hand.
func (a *App) release() {
	if !a.held {
		return
	}
	a.held = false
	if a.pending.ack != nil {
		close(a.pending.ack)
		a.pending.ack = nil
	}
}
