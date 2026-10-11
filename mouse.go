//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// mouse handles a real mouse or a trackpad — a Chromebook, a phone in
// desktop mode, a tablet with a keyboard case.
//
// Android reports which buttons are *held*, not which one changed, so the
// change has to be worked out by comparing with the last report. Getting
// this wrong produces a button that sticks down, which looks like a bug in
// the app rather than in the translation.
func (b *androidWindow) mouse(win *Window, e *ndk.InputEvent) bool {
	x, y := int(e.X(0)), int(e.Y(0))
	switch e.MotionAction() {
	case ndk.MotionScroll:
		// A wheel step is 1.0 per notch, and positive is away from the user,
		// which is the same direction antui calls up.
		if v := e.Axis(ndk.AxisVScroll, 0); v != 0 {
			win.Push(Event{Type: EventMouseWheel, Wheel: int(v), X: x, Y: y})
		}
		return true

	case ndk.MotionMove, ndk.MotionHoverMove:
		win.Push(Event{Type: EventMouseMove, X: x, Y: y, Mods: modsOf(e.Meta())})
		return true

	case ndk.MotionDown, ndk.MotionUp, ndk.MotionButtonPress, ndk.MotionButtonUp:
		now := e.Buttons()
		changed := now ^ b.buttons
		b.buttons = now
		for _, m := range []struct {
			bit    int32
			button MouseButton
		}{
			{ndk.ButtonPrimary, MouseLeft},
			{ndk.ButtonSecondary, MouseRight},
			{ndk.ButtonTertiary, MouseMiddle},
		} {
			if changed&m.bit == 0 {
				continue
			}
			kind := EventMouseUp
			if now&m.bit != 0 {
				kind = EventMouseDown
			}
			win.Push(Event{Type: kind, Button: m.button, X: x, Y: y, Mods: modsOf(e.Meta())})
		}
		return true
	}
	return false
}
