package backend

// EventType says which kind of thing happened. The fields of Event that
// carry meaning depend on it; the rest are zero.
type EventType int

// The event kinds. The comment on each names the fields it fills in.
const (
	EventNone       EventType = iota
	EventClose                // the user asked to close the window
	EventResize               // Width, Height
	EventKeyDown              // Key, Mods, Repeat
	EventKeyUp                // Key, Mods
	EventText                 // Text, Rune
	EventMouseDown            // Button, X, Y, Mods
	EventMouseUp              // Button, X, Y, Mods
	EventMouseMove            // X, Y, DX, DY
	EventMouseWheel           // Wheel (positive = upwards)
	EventFocus                // Focused
	EventExpose               // the window needs to be redrawn
	EventDropFiles            // Files, X, Y: files dragged onto the window

	// The touch events. Each carries TouchID, X, Y, Pressure, TouchSize and
	// Tool. A device with no touchscreen never produces them.
	EventTouchDown
	EventTouchMove
	EventTouchUp
	EventTouchCancel // the system took the touch away; it was not a tap
)

// Event is one thing the window system reported.
type Event struct {
	Type    EventType
	Key     Key
	Button  MouseButton
	Mods    Mod
	Repeat  bool // the key is auto-repeating rather than newly pressed
	X, Y    int  // mouse position
	DX, DY  int  // mouse movement since the last move event
	Wheel   int  // wheel steps
	Width   int  // new size, on EventResize
	Height  int
	Focused bool   // gained focus, on EventFocus
	Rune    rune   // character typed, on EventText
	Text    string // the same character in UTF-8

	// Files dragged onto the window, on EventDropFiles: paths on this
	// machine, in the order the system gave them.
	Files []string

	// TouchID identifies the finger, on the touch events. It stays with that
	// finger until it lifts.
	TouchID int
	// Pressure and TouchSize are 0 to 1, and 0 where the device does not say.
	Pressure  float32
	TouchSize float32
	// Tool is what is touching: a finger, a stylus, a mouse.
	Tool Tool
}