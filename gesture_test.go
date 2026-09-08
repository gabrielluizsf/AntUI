package antui

import (
	"math"
	"testing"
)

// frame plays one frame: the per-frame reset, then the events the backend
// would have pushed, then the recogniser. It is what Begin does, with the
// clock handed in.
func frame(win *Window, now float64, events ...Event) Gestures {
	win.beginTouchFrame()
	win.queue = win.queue[:0]
	for _, e := range events {
		win.Push(e)
	}
	win.recognize(now)
	return win.Gestures()
}

func down(id, x, y int) Event {
	return Event{Type: EventTouchDown, TouchID: id, X: x, Y: y, Tool: ToolFinger}
}
func move(id, x, y int) Event {
	return Event{Type: EventTouchMove, TouchID: id, X: x, Y: y, Tool: ToolFinger}
}
func up(id, x, y int) Event {
	return Event{Type: EventTouchUp, TouchID: id, X: x, Y: y, Tool: ToolFinger}
}

func TestGestureTap(t *testing.T) {
	win := &Window{}
	if g := frame(win, 0, down(1, 10, 10)); g.Tap {
		t.Error("a tap was reported before the finger lifted")
	}
	g := frame(win, 0.1, up(1, 10, 10))
	if !g.Tap || g.Count != 1 {
		t.Fatalf("lifting after 100ms gave Tap=%v Count=%d, want a single tap", g.Tap, g.Count)
	}
	if g.X != 10 || g.Y != 10 {
		t.Errorf("the tap is at (%d, %d), want (10, 10)", g.X, g.Y)
	}
	// And it is gone on the next frame, like every other one-frame state.
	if frame(win, 0.2).Tap {
		t.Error("the tap was still reported a frame later")
	}
}

func TestGestureDoubleTap(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	frame(win, 0.05, up(1, 10, 10))
	frame(win, 0.15, down(2, 12, 11))
	if g := frame(win, 0.2, up(2, 12, 11)); g.Count != 2 {
		t.Fatalf("two quick taps gave Count=%d, want 2", g.Count)
	}
	// A third, too late to join them, starts again rather than making three.
	frame(win, 2.0, down(3, 12, 11))
	if g := frame(win, 2.05, up(3, 12, 11)); g.Count != 1 {
		t.Fatalf("a tap two seconds later gave Count=%d, want 1", g.Count)
	}
}

func TestGestureDoubleTapMustBeNearby(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	frame(win, 0.05, up(1, 10, 10))
	frame(win, 0.1, down(2, 300, 400))
	if g := frame(win, 0.15, up(2, 300, 400)); g.Count != 1 {
		t.Fatalf("a second tap across the screen gave Count=%d, want 1", g.Count)
	}
}

func TestGestureLongPress(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	if g := frame(win, 0.3); g.LongPress {
		t.Error("a long press fired after 300ms")
	}
	if g := frame(win, 0.6); !g.LongPress {
		t.Fatal("no long press after 600ms")
	}
	// Once, not once a frame.
	if g := frame(win, 0.7); g.LongPress {
		t.Error("the long press fired again on the next frame")
	}
	// And what follows is not also a tap.
	if g := frame(win, 0.8, up(1, 10, 10)); g.Tap {
		t.Error("lifting after a long press was also reported as a tap")
	}
}

func TestGestureDragIsNotATap(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	g := frame(win, 0.05, move(1, 90, 10))
	if !g.Drag || g.DX != 80 {
		t.Fatalf("moving 80 pixels gave Drag=%v DX=%d", g.Drag, g.DX)
	}
	if g := frame(win, 0.1, up(1, 90, 10)); g.Tap {
		t.Error("a finger that moved 80 pixels was reported as a tap")
	}
}

func TestGestureSmallMoveIsStillATap(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	// Inside the slop: a finger never lands and lifts on exactly one pixel.
	frame(win, 0.05, move(1, 13, 12))
	if g := frame(win, 0.1, up(1, 13, 12)); !g.Tap {
		t.Fatal("a finger that wandered three pixels was not a tap")
	}
}

func TestGestureFling(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	frame(win, 0.1, move(1, 100, 10))
	g := frame(win, 0.15, up(1, 100, 10))
	if !g.Fling {
		t.Fatal("letting go at speed was not a fling")
	}
	if g.VX <= 0 {
		t.Errorf("the fling is going at %.0f pixels a second to the right, want positive", g.VX)
	}
}

func TestGestureSlowReleaseIsNotAFling(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	frame(win, 0.5, move(1, 40, 10))
	// Held still for a while before letting go: the velocity decays and
	// there is nothing to throw.
	frame(win, 1.0, move(1, 40, 10))
	frame(win, 1.5, move(1, 40, 10))
	if g := frame(win, 2.0, up(1, 40, 10)); g.Fling {
		t.Errorf("a finger held still before lifting flung at %.0f, %.0f", g.VX, g.VY)
	}
}

func TestGestureCancelledTouchIsNothing(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 10, 10))
	g := frame(win, 0.1, Event{Type: EventTouchCancel, TouchID: 1, X: 10, Y: 10})
	if g.Tap || g.Fling {
		t.Errorf("a cancelled touch reported Tap=%v Fling=%v; the system took it "+
			"away and the user did nothing", g.Tap, g.Fling)
	}
}

func TestGesturePinch(t *testing.T) {
	win := &Window{}
	// The first frame with two fingers only records where they are: there is
	// no previous distance to compare against yet.
	if g := frame(win, 0, down(1, 0, 0), down(2, 100, 0)); g.Pinch {
		t.Fatalf("a pinch was reported on the first frame, scale %.2f", g.Scale)
	}
	g := frame(win, 0.05, move(2, 200, 0))
	if !g.Pinch {
		t.Fatal("moving the fingers twice as far apart was not a pinch")
	}
	if math.Abs(g.Scale-2) > 0.001 {
		t.Errorf("scale %.3f, want 2", g.Scale)
	}
	if g.X != 100 || g.Y != 0 {
		t.Errorf("the pinch is at (%d, %d), want the point between the fingers, (100, 0)", g.X, g.Y)
	}
}

func TestGestureRotate(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 0, 0), down(2, 100, 0))
	g := frame(win, 0.05, move(2, 0, 100))
	if !g.Rotate {
		t.Fatal("turning the fingers a quarter turn was not a rotation")
	}
	if math.Abs(g.Angle-math.Pi/2) > 0.001 {
		t.Errorf("angle %.4f, want %.4f", g.Angle, math.Pi/2)
	}
}

// Crossing from just under pi to just over -pi is a small turn, not a full
// one in the other direction. This is the bug every rotation recogniser has
// until it is written down.
func TestGestureRotateWrapsTheShortWay(t *testing.T) {
	win := &Window{}
	frame(win, 0, down(1, 0, 0), down(2, -1000, 20))
	g := frame(win, 0.05, move(2, -1000, -20))
	if !g.Rotate {
		t.Fatal("no rotation across the wrap")
	}
	if math.Abs(g.Angle) > 0.1 {
		t.Errorf("a small turn across pi came out as %.4f radians", g.Angle)
	}
}

func TestGestureNothingOnADesktop(t *testing.T) {
	win := &Window{}
	// A mouse produces no touches, so there is nothing to recognise.
	g := frame(win, 0,
		Event{Type: EventMouseDown, Button: MouseLeft, X: 10, Y: 10},
		Event{Type: EventMouseUp, Button: MouseLeft, X: 10, Y: 10})
	if g != (Gestures{}) {
		t.Errorf("a mouse click produced gestures: %+v", g)
	}
}