//go:build android

package antui

// pump takes every event that is waiting, and then — if there is nothing to
// draw on — waits for one.
//
// That wait is the whole of what makes an app behave when it is not in
// front. A frame loop with no surface has nothing to do, and spinning
// through it burns a battery and gets the process killed; parked on a
// channel it costs nothing until the app is resumed.
func (b *androidWindow) pump(win *Window) {
	if b.activity == nil {
		return
	}
	for {
		e, ok := b.activity.Poll()
		if !ok {
			break
		}
		b.handle(win, e)
	}
	b.drainInput(win)
	for b.surface == nil && !win.shouldClose {
		e, ok := b.activity.Next()
		if !ok {
			win.Push(Event{Type: EventClose})
			break
		}
		b.handle(win, e)
	}
	// Let go of the last event before the frame starts. Everything after
	// this belongs to the app, and the UI thread must not be blocked for the
	// length of it.
	b.activity.Release()
	b.sync(win)
}
