package antui

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// TestAFailedPresentRetriesAFullRedrawNextFrame guards the bug behind a
// window that goes blank partway through a big resize: the platform can drop
// a frame (a connection hiccup, a dropped buffer) without the drawing being
// wrong at all. If the next frame trusted the diff against what was measured
// rather than what actually reached the screen, a frame that redraws nothing
// new would see no change and send nothing — leaving the display stuck
// showing whatever the failed present left behind, forever.
func TestAFailedPresentRetriesAFullRedrawNextFrame(t *testing.T) {
	win, stub := newTestWindow(t, 64, 64)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	if !win.Begin() {
		t.Fatal("window closed")
	}
	drawTestFrame(win)
	if !win.End() {
		t.Fatal("the first frame presented nothing, but the screen was empty")
	}

	// The platform refuses this frame's present, the way a dropped X11
	// connection would partway through a big resize's full redraw.
	stub.refusePresent = true
	if !win.Begin() {
		t.Fatal("window closed")
	}
	win.FillRect(canvas.Area{X: 10, Y: 10, Width: 6, Height: 6}, canvas.White)
	if win.End() {
		t.Fatal("End reported success for a present the platform refused")
	}

	// The platform recovers, but this frame draws nothing new on top of the
	// one that failed to arrive.
	stub.refusePresent = false
	if !win.Begin() {
		t.Fatal("window closed")
	}
	if !win.End() {
		t.Error("a frame after a dropped present sent nothing: the screen is stuck with whatever the failed present left behind")
	}
}
