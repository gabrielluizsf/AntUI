package antui

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// TestTheLoopWaitsForItsOwnFrameRate is the half of "does not spin" that lives
// in the window: a program that asked for 60 frames a second gets 60 frames a
// second, and the wait between them is what makes it so. Before the loop had a
// budget it went round Begin/End as fast as the drawing took, which on a fast
// frame is several hundred times a second and is heard as a fan.
func TestTheLoopWaitsForItsOwnFrameRate(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.SetFPS(60)

	start := time.Now()
	for range 10 {
		win.pace()
	}
	elapsed := time.Since(start)

	if want := 9 * time.Second / 60; elapsed < want*3/4 || elapsed > want*5/4 {
		t.Errorf("ten frames at 60 Hz took %v, want about %v", elapsed, want)
	}
}

// TestTheLoopDoesNotRepayAnOverrun is what a deadline grid buys over sleeping a
// fixed time after each frame: a frame that took too long has already missed
// its slot, and waiting the remainder of it would put every frame after it
// later and later, for as long as the app stayed slow.
func TestTheLoopDoesNotRepayAnOverrun(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.SetFPS(60)
	win.sleep = func(time.Duration) {}

	// Ten frames, each of which took longer than the frame time.
	start := time.Now()
	for range 10 {
		win.pace()
		time.Sleep(20 * time.Millisecond)
	}
	if elapsed := time.Since(start); elapsed < 190*time.Millisecond {
		t.Errorf("ten frames that each overran their slot finished in %v: the loop repaid the overrun on top of the work", elapsed)
	}
}

// TestAStillFrameIsFree is the point of measuring a frame against the one on
// the screen: a window nobody is touching neither sends anything nor collects
// anything. The first frame goes out — the screen is empty — and every frame
// after it asks the platform for nothing at all.
func TestAStillFrameIsFree(t *testing.T) {
	win, stub := newTestWindow(t, 400, 300)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	frame := func() bool {
		if !win.Begin() {
			t.Fatal("window closed")
		}
		drawTestFrame(win)
		return win.End()
	}

	if !frame() {
		t.Fatal("the first frame presented nothing, but the screen was empty")
	}
	allocs := testing.AllocsPerRun(20, func() {
		if frame() {
			t.Error("a frame that draws the picture the screen already has was presented")
		}
	})
	if allocs > 0 {
		t.Errorf("a still frame allocates %v times; a frame that presents nothing should cost nothing at all", allocs)
	}
	if len(stub.presented) != 1 {
		t.Errorf("the platform was asked to present %d times, want the first frame only", len(stub.presented))
	}
}

// TestAFrameIsPresentedWhenSomethingMoves is the other half: measuring a frame
// against the one on the screen is only worth anything if a frame that differs
// from it gets through.
func TestAFrameIsPresentedWhenSomethingMoves(t *testing.T) {
	win, stub := newTestWindow(t, 400, 300)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	for i := range 3 {
		if !win.Begin() {
			t.Fatal("window closed")
		}
		drawTestFrame(win)
		win.FillRect(10+i, 10, 6, 6, canvas.White)
		if !win.End() {
			t.Fatalf("frame %d moved a pixel and presented nothing", i)
		}
	}
	if len(stub.presented) != 3 {
		t.Errorf("the platform was asked to present %d times, want one per frame that moved", len(stub.presented))
	}
}
