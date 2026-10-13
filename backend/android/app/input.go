//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// Input takes everything waiting in the input queue and hands each event to
// handle, which reports whether the app dealt with it. An event the app does
// not deal with goes back to the system — which is how the back button still
// closes an app that ignored it.
//
// It must be called from the same goroutine as [App.Next], which is the one
// the looper belongs to.
func (a *App) Input(handle func(*ndk.InputEvent) bool) {
	if a.queue == nil || a.loop == nil {
		return
	}
	// Give the looper a turn so it notices the queue is readable. Emptying
	// it is the loop below; this only wakes the machinery.
	a.loop.Poll(0)
	for {
		e, ok := a.queue.Next()
		if !ok {
			return
		}
		// The system gets first refusal on everything. This is how the
		// on-screen keyboard sees a key at all, and skipping it produces an
		// app that works until somebody types.
		if a.queue.PreDispatch(e) {
			continue
		}
		a.queue.Finish(e, handle(e))
	}
}
