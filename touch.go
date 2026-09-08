package antui

import "github.com/gabrielluizsf/antui/backend"

// Tool is what is touching the screen.
type Tool = backend.Tool

// The tools. A device that does not say reports ToolUnknown rather than
// guessing at a finger.
const (
	ToolUnknown = backend.ToolUnknown
	ToolFinger  = backend.ToolFinger
	ToolStylus  = backend.ToolStylus
	ToolMouse   = backend.ToolMouse
	ToolEraser  = backend.ToolEraser // the far end of a stylus, used to rub out
)

// Touch is one finger on the screen, from the moment it lands to the moment
// it lifts.
//
// It is not a mouse and is not reported as one. A mouse has one position and
// a set of buttons; fingers arrive several at a time, each with its own
// history, and a program that wants to pinch or to let two people play on
// one screen needs to tell them apart. For everything simpler, the first
// finger down is also reported as the left mouse button, so a program
// written for a desktop works on a phone without knowing.
type Touch struct {
	// ID stays with this finger for as long as it is down, and may be reused
	// by a later one. It is not an index: a finger's position in
	// [Window.Touches] changes as others come and go.
	ID int

	// X and Y are where it is now, in canvas pixels.
	X, Y int
	// DX and DY are how far it moved this frame.
	DX, DY int
	// StartX and StartY are where it landed, which is what a swipe is
	// measured from.
	StartX, StartY int

	// Pressure is 0 to 1 on hardware that measures it. On hardware that does
	// not it is 1 the whole time the finger is down, so it cannot be used to
	// find out whether the device reports pressure at all.
	Pressure float32
	// Size is how much of the screen the touch covers, 0 to 1. It is how a
	// thumb is told from a fingertip.
	Size float32
	// Tool is what is doing the touching.
	Tool Tool

	// Began is set for the one frame the finger landed on, and Ended for the
	// one frame it lifted on. A touch is still in [Window.Touches] on the
	// frame it ended, so nothing has to be caught between frames.
	Began bool
	Ended bool
	// Cancelled says the system took the touch away rather than the user
	// lifting it — a notification pulled down over the app, a gesture the
	// system claimed. A tap must not be counted from a cancelled touch.
	Cancelled bool
}

// Touches is every finger on the screen this frame, in the order they
// landed. The slice is the window's own and is rewritten every frame, so
// anything kept has to be copied.
func (win *Window) Touches() []Touch { return win.touches }

// TouchCount is how many fingers are on the screen.
func (win *Window) TouchCount() int { return len(win.touches) }

// TouchScreen reports whether this window is being driven by a touchscreen,
// which is what decides whether a game has to put controls on the screen.
//
// It is not a question about the hardware, it is a question about the
// player: a laptop with a touchscreen nobody uses should be given a keyboard
// game, and a tablet with a keyboard attached still has a screen to press.
// So the answer is "a finger has touched this window", which is false until
// one does — and on a platform where there is nothing else, true from the
// start.
//
// A game that wants to decide for itself should ignore this and set its own
// flag; this is the default, not the rule.
func (win *Window) TouchScreen() bool {
	if win == nil {
		return false
	}
	return win.touchFirst || win.sawTouch
}

// setTouchFirst is how a backend says there is no other kind of input on
// this platform, so that the first frame already knows rather than waiting
// for a finger.
func (win *Window) SetTouchFirst() { win.touchFirst = true }

// TouchByID finds a touch by the identifier that stays with it.
func (win *Window) TouchByID(id int) (Touch, bool) {
	for _, t := range win.touches {
		if t.ID == id {
			return t, true
		}
	}
	return Touch{}, false
}

// beginTouchFrame is the per-frame reset: the fingers that lifted last frame
// are gone, and nothing has moved yet.
func (win *Window) beginTouchFrame() {
	kept := win.touches[:0]
	for _, t := range win.touches {
		if t.Ended {
			continue
		}
		t.Began = false
		t.DX, t.DY = 0, 0
		kept = append(kept, t)
	}
	win.touches = kept
}

// touchIndex finds a live touch, or -1.
func (win *Window) touchIndex(id int) int {
	for i := range win.touches {
		if win.touches[i].ID == id {
			return i
		}
	}
	return -1
}

// pushTouch folds a touch event into the window's touch state. It is called
// from push, so the state is right whether or not the queue overflowed.
func (win *Window) pushTouch(ev Event) {
	win.sawTouch = true
	i := win.touchIndex(ev.TouchID)
	switch ev.Type {
	case EventTouchDown:
		// A second down for a finger already on the screen is the system
		// repeating itself, which happens. Taking it as a move keeps the
		// history rather than starting a new one.
		if i >= 0 {
			break
		}
		win.touches = append(win.touches, Touch{
			ID:       ev.TouchID,
			X:        ev.X,
			Y:        ev.Y,
			StartX:   ev.X,
			StartY:   ev.Y,
			Pressure: ev.Pressure,
			Size:     ev.TouchSize,
			Tool:     ev.Tool,
			Began:    true,
		})
		return
	case EventTouchUp, EventTouchCancel:
		if i < 0 {
			return
		}
		win.touches[i].Ended = true
		win.touches[i].Cancelled = ev.Type == EventTouchCancel
	}
	if i < 0 {
		// A move for a finger that was never reported down. It happens when
		// the app gains a surface mid-gesture; treating it as a landing is
		// better than dropping the finger on the floor.
		win.touches = append(win.touches, Touch{
			ID: ev.TouchID, StartX: ev.X, StartY: ev.Y, Began: true,
			Tool: ev.Tool,
		})
		i = len(win.touches) - 1
	}
	t := &win.touches[i]
	t.DX += ev.X - t.X
	t.DY += ev.Y - t.Y
	t.X, t.Y = ev.X, ev.Y
	if ev.Pressure > 0 {
		t.Pressure = ev.Pressure
	}
	if ev.TouchSize > 0 {
		t.Size = ev.TouchSize
	}
}
