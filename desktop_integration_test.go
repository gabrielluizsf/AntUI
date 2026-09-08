//go:build linux || darwin || windows

package antui

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// TestDesktopWindowRunsAnApp is the integration test: it opens a real window
// on the machine the suite is running on and drives it the way an app would,
// so the backends that cannot be exercised from a cross-compile — the X11
// socket, the Win32 message loop, the Cocoa run loop — are heard from by the
// only thing that can make them say anything, which is running them.
//
// It skips where there is no desktop to open a window on (a headless box, a
// container), and this build tag keeps it off the platforms with no window
// at all. What runs on the CI matrix is this file on each of the three.
//
// The assertions are strict where the library owns the whole path — open,
// draw, present, title, icon, clipboard — and merely report where a window
// manager has a vote, which is fullscreen and the exact delivered size.
func TestDesktopWindowRunsAnApp(t *testing.T) {
	win, err := OpenWith(Options{
		Title:  "AntUI integration test",
		Width:  640,
		Height: 480,
		Limits: Limits{MaxWidth: 1600, MaxHeight: 1200},
	})
	if err != nil {
		t.Skipf("no desktop to open a window on: %v", err)
	}
	defer win.Close()

	if win.Width() != 640 || win.Height() != 480 {
		t.Errorf("the window opened at %dx%d, want 640x480 (a smaller display shrinks it)",
			win.Width(), win.Height())
	}

	win.SetTitle("AntUI integration test")

	// An icon is part of opening a window on a desktop: the task bar and the
	// title bar both ask for one.
	big, err := canvas.NewCanvas(32, 32)
	if err != nil {
		t.Fatal(err)
	}
	big.Clear(canvas.RGB(0x3E, 0x63, 0xDD))
	small, err := canvas.NewCanvas(16, 16)
	if err != nil {
		t.Fatal(err)
	}
	small.Clear(canvas.RGB(0x3E, 0x63, 0xDD))
	if !win.SetIcon(big, small) {
		t.Error("the window system refused the icon")
	}

	// The frame loop, running the way an app runs it. Most of it is ordinary
	// drawing; a few frames in the window also goes fullscreen, and a few
	// later it is asked to resize.
	swatch := canvas.RGB(0x3E, 0x63, 0xDD)
	frames, events, drew := 0, 0, false
	resizeSeen := false
	deadline := time.Now().Add(4 * time.Second)
	for win.Begin() && frames < 240 {
		win.Clear(canvas.RGB(0x10, 0x12, 0x18))
		win.FillRect(20, 20, 200, 120, swatch)
		win.Circle(win.Width()/2, win.Height()/2, 60, canvas.White)
		win.Line(0, win.Height()-1, win.Width()-1, 0, canvas.RGBA(0xFF, 0, 0, 0xFF))
		win.Rect(4, 4, win.Width()-8, win.Height()-8, canvas.RGBA(0xFF, 0xFF, 0xFF, 0x80))
		win.Text(20, 160, "AntUI builds desktop apps the simple way.", canvas.White)
		win.End()
		frames++

		if win.Canvas().At(25, 25) == swatch {
			drew = true
		}
		win.Events(func(ev Event) bool {
			events++
			if ev.Type == EventResize {
				resizeSeen = true
			}
			return true
		})

		switch frames {
		case 15:
			if !win.SetFullscreen(true) {
				t.Error("fullscreen could not even be asked for")
			}
			if !win.Fullscreen() {
				t.Error("the window does not believe it asked for full screen")
			}
		case 20:
			win.SetFullscreen(false)
			win.SetSize(800, 600)
		}
		if time.Now().After(deadline) {
			break
		}
	}

	if frames == 0 {
		t.Fatal("the frame loop never ran a single frame")
	}
	if !drew {
		t.Error("a painted rectangle never made it into the framebuffer")
	}
	t.Logf("ran %d frames and heard %d events", frames, events)
	if !resizeSeen && win.Width() != 800 {
		t.Errorf("after SetSize(800, 600) the canvas is %dx%d and no resize arrived",
			win.Width(), win.Height())
	}
	if win.Fullscreen() {
		t.Error("the window still believes it is full screen")
	}

	// What the display said about itself, for the log rather than for a
	// verdict: numbers like these depend on where the machine is, and a
	// rectangle that does not report them is not wrong.
	if w, h, ok := win.DisplaySize(); ok {
		t.Logf("the display is %dx%d", w, h)
	} else {
		t.Log("the display did not say its size")
	}
	if hz := win.DisplayRefresh(); hz > 0 {
		t.Logf("the display refreshes at %d Hz", hz)
	}
	if scale := win.ContentScale(); scale > 1 {
		t.Logf("the display is scaled at %g", scale)
	}

	// The clipboard is a round trip: what the window puts there, it should
	// find again — the platforms tell it apart only by how the two halves
	// are asked, so a port that gets one ask right and the other wrong
	// surfaces exactly here.
	const clip = "AntUI clipboard round-trip"
	if !win.SetClipboardText(clip) {
		t.Error("the clipboard refused the text")
	} else if got := win.ClipboardText(); got != clip {
		t.Errorf("the clipboard came back %q, want %q", got, clip)
	}

	// The size lock, which every desktop enforces differently.
	win.SetFixedSize(640, 480)
	if win.Resizable() {
		t.Error("a fixed window still calls itself resizable")
	}
	win.SetLimits(Limits{})
	if !win.Resizable() {
		t.Error("removing the limits left the window locked")
	}

	// The platform's own font is found by asking the system, which is a
	// different mechanism on each desktop. A machine with nothing to answer
	// is not a failure — every program carries on with the built-in font —
	// but the log says which it was.
	if _, err := SystemFace(14); err != nil {
		t.Logf("no system font to draw with: %v", err)
	} else {
		t.Log("read the platform's own font")
	}
}