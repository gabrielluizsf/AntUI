package antui

import "github.com/gabrielluizsf/antui/backend"

// EventType says which kind of thing happened. The fields of Event that
// carry meaning depend on it; the rest are zero.
type EventType = backend.EventType

// The event kinds. The comment on each names the fields it fills in.
const (
	EventNone       = backend.EventNone
	EventClose      = backend.EventClose
	EventResize     = backend.EventResize
	EventKeyDown    = backend.EventKeyDown
	EventKeyUp      = backend.EventKeyUp
	EventText       = backend.EventText
	EventMouseDown  = backend.EventMouseDown
	EventMouseUp    = backend.EventMouseUp
	EventMouseMove  = backend.EventMouseMove
	EventMouseWheel = backend.EventMouseWheel
	EventFocus      = backend.EventFocus
	EventExpose     = backend.EventExpose
	EventDropFiles  = backend.EventDropFiles
	EventTouchDown    = backend.EventTouchDown
	EventTouchMove    = backend.EventTouchMove
	EventTouchUp      = backend.EventTouchUp
	EventTouchCancel  = backend.EventTouchCancel
)

// Event is one thing the window system reported. Reading events is optional:
// Begin has already folded them into the input state that Window.KeyDown and
// friends report, and most programs only need that.
type Event = backend.Event

// maxEvents bounds one frame's event queue. A frame that somehow produces
// more than this drops the excess rather than growing without limit — the
// input state is still correct, because it is updated as events arrive.
const maxEvents = 256