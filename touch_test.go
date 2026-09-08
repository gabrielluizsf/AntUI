package antui

import (
	"testing"
)

// The touch path keeps one record per finger: where it landed, where it is
// now, and whether it is still on the screen. Fingers arrive one event at a
// time, so the window has to fold them back together again.

// A finger's record has to survive from the frame it lands to the frame it
// lifts, because the code that reads touches gets one frame at a time.
func TestTouchLifecycle(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)

	win.pushTouch(Event{Type: EventTouchDown, TouchID: 7, X: 10, Y: 20,
		Pressure: 0.5, Tool: ToolFinger})
	if win.TouchCount() != 1 {
		t.Fatalf("after a down: %d touches", win.TouchCount())
	}
	touch, ok := win.TouchByID(7)
	if !ok {
		t.Fatal("the finger that landed is not findable")
	}
	if !touch.Began || touch.Pressure != 0.5 || touch.Tool != ToolFinger {
		t.Errorf("the down did not register: %+v", touch)
	}
	if touch.X != 10 || touch.Y != 20 || touch.StartX != 10 || touch.StartY != 20 {
		t.Errorf("the record does not hold where it landed: %+v", touch)
	}

	win.pushTouch(Event{Type: EventTouchMove, TouchID: 7, X: 30, Y: 40})
	win.pushTouch(Event{Type: EventTouchMove, TouchID: 7, X: 45, Y: 40})
	touch, _ = win.TouchByID(7)
	if touch.X != 45 || touch.Y != 40 || touch.DX != 35 || touch.DY != 20 {
		t.Errorf("moves do not add up: %+v", touch)
	}

	win.pushTouch(Event{Type: EventTouchUp, TouchID: 7, X: 45, Y: 40})
	// Ended touches are still visible on the frame they lifted, so the game
	// can read the last position.
	if touch, ok = win.TouchByID(7); !ok || !touch.Ended {
		t.Fatalf("the lift did not register: %+v", touch)
	}

	win.beginTouchFrame()
	if win.TouchCount() != 0 {
		t.Errorf("a lifted finger is still here next frame: %+v", win.Touches())
	}
}

// A new frame means a new reading: what was movement last frame is the
// finger's position now, and a tap cannot go on counting as a landing.
func TestTouchFrameBegin(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.pushTouch(Event{Type: EventTouchDown, TouchID: 1, X: 0, Y: 0})
	win.pushTouch(Event{Type: EventTouchMove, TouchID: 1, X: 5, Y: 5})

	win.beginTouchFrame()
	touch, _ := win.TouchByID(1)
	if touch.Began || touch.DX != 0 || touch.DY != 0 {
		t.Errorf("a new frame kept the last frame's movement: %+v", touch)
	}
	if win.TouchCount() != 1 {
		t.Fatalf("a live finger vanished: %d", win.TouchCount())
	}
}

// A second down for a finger already on the screen is the system repeating
// itself, which happens; the window takes it as a move rather than starting
// a second finger on top of itself.
func TestDuplicateDownKeepsOneFinger(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.pushTouch(Event{Type: EventTouchDown, TouchID: 1, X: 0, Y: 0})
	win.pushTouch(Event{Type: EventTouchDown, TouchID: 1, X: 8, Y: 9})
	if win.TouchCount() != 1 {
		t.Fatalf("a repeated down made %d fingers", win.TouchCount())
	}
	touch, _ := win.TouchByID(1)
	if touch.X != 8 || touch.Y != 9 {
		t.Errorf("the repeated down was not folded into a move: %+v", touch)
	}
}

// A move for a finger that was never reported down happens when the window
// gains a surface mid-gesture; treating it as a landing is better than
// dropping the finger on the floor.
func TestMoveForANewFingerLandsIt(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.pushTouch(Event{Type: EventTouchMove, TouchID: 3, X: 12, Y: 14})
	touch, ok := win.TouchByID(3)
	if !ok || !touch.Began || touch.StartX != 12 || touch.StartY != 14 {
		t.Errorf("a move is not a landing when there is no finger: %+v", touch)
	}
}

// A cancelled touch is not a tap: the system took it away rather than the
// user lifting it, and a game counting taps from it would count a
// notification pulled down over the app.
func TestCancelledTouchesAreNotTaps(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.pushTouch(Event{Type: EventTouchDown, TouchID: 4, X: 0, Y: 0})
	win.pushTouch(Event{Type: EventTouchCancel, TouchID: 4})
	touch, ok := win.TouchByID(4)
	if !ok || !touch.Ended || !touch.Cancelled {
		t.Fatalf("the cancel did not register: %+v", touch)
	}
}

// TouchScreen is learned, not asked: it is false until a finger has touched
// the window, or the platform says there is no other kind of input.
func TestTouchScreenIsLearned(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	if win.TouchScreen() {
		t.Fatal("a window that has never seen a finger is a touch screen")
	}
	win.pushTouch(Event{Type: EventTouchDown, TouchID: 1, X: 0, Y: 0})
	if !win.TouchScreen() {
		t.Fatal("a window that has just seen a finger is not a touch screen")
	}

	win2, _ := newTestWindow(t, 320, 200)
	win2.SetTouchFirst()
	if !win2.TouchScreen() {
		t.Fatal("SetTouchFirst did not say what the platform knows")
	}
}