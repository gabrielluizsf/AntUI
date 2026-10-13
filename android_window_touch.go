//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// touch turns one motion report into touch events, and drives the mouse from
// the first finger down.
//
// The action says what happened and, for the two "pointer" variants, which
// finger it happened to — as an **index into this event**, not as an id. A
// move reports every finger at once; everything else reports one.
func (b *androidWindow) touch(win *Window, e *ndk.InputEvent) bool {
	switch e.MotionAction() {
	case ndk.MotionDown, ndk.MotionPointerDown:
		b.touchAt(win, e, e.MotionIndex(), EventTouchDown)
	case ndk.MotionUp, ndk.MotionPointerUp:
		b.touchAt(win, e, e.MotionIndex(), EventTouchUp)
	case ndk.MotionMove:
		for i := range e.PointerCount() {
			b.touchAt(win, e, i, EventTouchMove)
		}
	case ndk.MotionCancel:
		for i := range e.PointerCount() {
			b.touchAt(win, e, i, EventTouchCancel)
		}
	default:
		return false
	}
	return true
}
