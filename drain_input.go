//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// drainInput takes everything the system has and turns it into events the
// window core understands. It runs once a frame, from pump.
func (b *androidWindow) drainInput(win *Window) {
	if b.activity == nil {
		return
	}
	b.activity.Input(func(e *ndk.InputEvent) bool {
		switch e.Kind() {
		case ndk.EventKey:
			return b.key(win, e)
		case ndk.EventMotion:
			return b.motion(win, e)
		}
		return false
	})
}
