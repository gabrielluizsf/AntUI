//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// touchAt pushes one finger's event, and the mouse event that shadows it.
func (b *androidWindow) touchAt(win *Window, e *ndk.InputEvent, index int, kind EventType) {
	if index < 0 || index >= e.PointerCount() {
		return
	}
	id := e.PointerID(index)
	x, y := int(e.X(index)), int(e.Y(index))

	win.Push(Event{
		Type:      kind,
		TouchID:   id,
		X:         x,
		Y:         y,
		Pressure:  e.Pressure(index),
		TouchSize: e.Size(index),
		Tool:      toolOf(e.ToolType(index)),
	})

	// One finger — the first one down — is also reported as the left mouse
	// button, so that a program written for a desktop works on a phone with
	// nothing added. The rest are touches only: a mouse has one position,
	// and pretending otherwise would make two fingers look like one jumping
	// between them.
	switch kind {
	case EventTouchDown:
		if b.primary >= 0 {
			return
		}
		b.primary = id
		win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: x, Y: y})
	case EventTouchMove:
		if id != b.primary {
			return
		}
		win.Push(Event{Type: EventMouseMove, X: x, Y: y})
	case EventTouchUp, EventTouchCancel:
		if id != b.primary {
			return
		}
		b.primary = -1
		win.Push(Event{Type: EventMouseUp, Button: MouseLeft, X: x, Y: y})
	}
}
