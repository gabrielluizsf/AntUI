package antui

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// The clamp is the whole of the arithmetic and is where every platform's
// answer comes from, so it is tested on its own before any window exists.
func TestLimitsClamp(t *testing.T) {
	cases := []struct {
		name   string
		limits Limits
		w, h   int
		wantW  int
		wantH  int
	}{
		{"no limits leave a size alone", Limits{}, 800, 600, 800, 600},
		{"a minimum lifts a window that is too small",
			Limits{MinWidth: 640, MinHeight: 480}, 320, 200, 640, 480},
		{"a minimum leaves a window that is big enough",
			Limits{MinWidth: 640, MinHeight: 480}, 1280, 720, 1280, 720},
		{"a maximum brings a window back down",
			Limits{MaxWidth: 1280, MaxHeight: 720}, 1920, 1080, 1280, 720},
		{"a bound on one axis leaves the other free",
			Limits{MinWidth: 640}, 320, 200, 640, 200},
		// A maximum below a minimum is what a display imposes on a window
		// asking for more than fits, and the display is the one that decides.
		{"a maximum below a minimum wins",
			Limits{MinWidth: 1920, MaxWidth: 1366}, 1920, 768, 1366, 768},
		// Nothing is ever allowed to reach zero: a canvas of no pixels is not
		// a window, it is a crash waiting for the first draw.
		{"a size of nothing becomes a pixel", Limits{}, 0, -5, 1, 1},

		{"an aspect makes the height follow the width",
			Limits{Aspect: 16.0 / 9.0}, 1600, 400, 1600, 900},
		// 1366 is not a multiple of 16, so 16:9 lands on 768.375 and the
		// window is a third of a pixel from the shape it asked for. Nearest
		// is the right answer; truncating would drift a pixel every drag.
		{"an aspect rounds rather than truncating",
			Limits{Aspect: 16.0 / 9.0}, 1366, 999, 1366, 768},
		{"a ratio that would break a maximum bends instead",
			Limits{Aspect: 16.0 / 9.0, MaxHeight: 720}, 1600, 400, 1280, 720},
		{"a ratio that would break a minimum bends instead",
			Limits{Aspect: 16.0 / 9.0, MinHeight: 720}, 800, 200, 1280, 720},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w, h := c.limits.Clamp(c.w, c.h)
			if w != c.wantW || h != c.wantH {
				t.Errorf("%dx%d clamped to %dx%d, want %dx%d",
					c.w, c.h, w, h, c.wantW, c.wantH)
			}
		})
	}

	// Clamping is idempotent: a size already inside the limits comes back as
	// itself, which is what stops a resize loop where each pass nudges the
	// window and provokes the next one.
	limits := Limits{MinWidth: 640, MinHeight: 480, MaxWidth: 1920, Aspect: 4.0 / 3.0}
	w, h := limits.Clamp(1000, 100)
	if w2, h2 := limits.Clamp(w, h); w2 != w || h2 != h {
		t.Errorf("clamping twice gave %dx%d then %dx%d", w, h, w2, h2)
	}
}

func TestFitDisplay(t *testing.T) {
	// Asking for more than the screen has is an ordinary Tuesday, not an
	// error: a 1920x1080 game opened on a 1366x768 laptop.
	if w, h := fitDisplay(1920, 1080, 1366, 768); w != 1366 || h != 768 {
		t.Errorf("got %dx%d, want 1366x768", w, h)
	}
	if w, h := fitDisplay(800, 600, 1920, 1080); w != 800 || h != 600 {
		t.Errorf("a window that fits was changed to %dx%d", w, h)
	}
	// A display that says nothing is not a display of nothing.
	if w, h := fitDisplay(800, 600, 0, 0); w != 800 || h != 600 {
		t.Errorf("an unknown display shrank the window to %dx%d", w, h)
	}
}

func TestOpenWithFitsTheDisplay(t *testing.T) {
	win, stub := newTestWindow(t, 100, 100)
	stub.displayW, stub.displayH = 1366, 768

	// SetSize goes through the same reduction, which is the path a game takes
	// when it changes resolution at runtime.
	if !win.SetSize(1920, 1080) {
		t.Fatal("the request could not be made")
	}
	if len(stub.sized) != 1 {
		t.Fatalf("the backend was asked %d times", len(stub.sized))
	}
	if got := stub.sized[0]; got.Width != 1366 || got.Height != 768 {
		t.Errorf("asked the backend for %dx%d, want 1366x768", got.Width, got.Height)
	}

	// And the canvas is not touched by the request. It moves when the resize
	// actually arrives, which is what keeps a window manager that ignored the
	// request from leaving the game drawing at a size it is not.
	if win.Width() != 100 || win.Height() != 100 {
		t.Errorf("the canvas moved to %dx%d on a request alone", win.Width(), win.Height())
	}
}

func TestSetLimitsReachesTheBackend(t *testing.T) {
	win, stub := newTestWindow(t, 800, 600)

	win.SetLimits(Limits{MinWidth: 640, MinHeight: 480})
	if stub.limits.MinWidth != 640 || stub.limits.MinHeight != 480 {
		t.Errorf("the backend got %+v", stub.limits)
	}
	if !win.Resizable() {
		t.Error("a window with a minimum stopped being resizable")
	}

	// A window already outside the new bounds is brought back inside them,
	// which is what makes SetLimits usable after the window is open.
	stub.sized = nil
	win.SetLimits(Limits{MaxWidth: 640, MaxHeight: 480})
	if len(stub.sized) != 1 || stub.sized[0].Width != 640 {
		t.Errorf("the window was not brought inside its new limits: %v", stub.sized)
	}
}

// A fixed window's smallest and largest size are both the size it has, and
// the backend is told so rather than being left to work it out.
func TestFixedSize(t *testing.T) {
	win, stub := newTestWindow(t, 800, 600)
	win.SetFixedSize(1280, 720)

	if win.Resizable() {
		t.Error("a fixed window says it can be resized")
	}
	if !stub.limits.Fixed || !stub.limits.NoMaximize {
		t.Errorf("the backend got %+v", stub.limits)
	}
	if stub.limits.MinWidth != 1280 || stub.limits.MaxWidth != 1280 ||
		stub.limits.MinHeight != 720 || stub.limits.MaxHeight != 720 {
		t.Errorf("the lock did not reach the backend as one size: %+v", stub.limits)
	}

	// Fixed on its own — no numbers — locks the window at whatever it is now,
	// resolved when it is asked for rather than frozen when it was set.
	loose, _ := newTestWindow(t, 800, 600)
	loose.SetLimits(Limits{Fixed: true})
	if got := loose.Bounds(); got.MinWidth != 800 || got.MaxHeight != 600 {
		t.Errorf("a bare lock resolved to %+v", got)
	}
	loose.ResizeCanvas(1024, 768)
	if got := loose.Bounds(); got.MinWidth != 1024 || got.MaxHeight != 768 {
		t.Errorf("the lock did not follow the window: %+v", got)
	}
}

// The half of this phase that matters: a resize that arrives at a locked
// window anyway — because a tiling window manager sent one regardless — is
// accepted and drawn correctly. The lock is a request, and the machines where
// it does not hold are the ones this must not break.
func TestALockedWindowStillFollowsAResizeItDidNotAskFor(t *testing.T) {
	win, _ := newTestWindow(t, 640, 360)
	win.SetFixedSize(640, 360)

	// What i3 does: the window is placed at the size the layout says.
	if !win.ResizeCanvas(1200, 300) {
		t.Fatal("a locked window refused the size it was given")
	}
	if win.Width() != 1200 || win.Height() != 300 {
		t.Errorf("the window reports %dx%d, want the 1200x300 it was given",
			win.Width(), win.Height())
	}
	if win.Canvas().Width != 1200 || win.Canvas().Height != 300 {
		t.Errorf("the canvas is %dx%d and the window is %dx%d",
			win.Canvas().Width, win.Canvas().Height, win.Width(), win.Height())
	}

	// And drawing still lands inside it rather than off the end of the old
	// one, which is the failure this is really guarding against.
	win.Clear(canvas.White)
	win.Canvas().Put(1199, 299, canvas.Black)
	if got := win.Canvas().At(1199, 299); got != canvas.Black {
		t.Errorf("the far corner of the new size is %v", got)
	}

	// And the two numbers stay apart, which is the distinction this whole
	// phase turns on: the lock is still what the game asked for, and the size
	// is still what actually happened. A lock that quietly rewrote itself to
	// whatever the window manager did would be no lock at all — the next
	// window manager that reads hints would place the window at 1200x300.
	if got := win.Bounds(); got.MinWidth != 640 || got.MaxHeight != 360 {
		t.Errorf("the lock became %+v rather than staying what was asked for", got)
	}
}

func TestContentScale(t *testing.T) {
	win, stub := newTestWindow(t, 1280, 720)

	// A system that does not say answers 1, so arithmetic on it is always
	// safe and no caller needs a special case.
	if got := win.ContentScale(); got != 1 {
		t.Errorf("an unscaled display reports %v", got)
	}
	if w, h := win.SizeInPoints(); w != 1280 || h != 720 {
		t.Errorf("at 1x the window appears %dx%d", w, h)
	}

	win.scale = 2
	if w, h := win.SizeInPoints(); w != 640 || h != 360 {
		t.Errorf("at 2x a 1280x720 window appears %dx%d, want 640x360", w, h)
	}
	win.scale = 1.5
	if w, h := win.SizeInPoints(); w != 853 || h != 480 {
		t.Errorf("at 1.5x the window appears %dx%d", w, h)
	}
	_ = stub

	var none *Window
	if none.ContentScale() != 1 {
		t.Error("a window that is not there has no scale to divide by")
	}
}

// A size given in points is what makes a locked window the same size on every
// machine instead of a postage stamp on a 4K laptop.
func TestPointsScaleWithTheDisplay(t *testing.T) {
	if got := scaleTo(1280, 2); got != 2560 {
		t.Errorf("1280 points at 2x is %d", got)
	}
	if got := scaleTo(1280, 1); got != 1280 {
		t.Errorf("an unscaled display changed 1280 to %d", got)
	}
	if got := scaleTo(1280, 0); got != 1280 {
		t.Errorf("a display that says nothing changed 1280 to %d", got)
	}
	if got := scaleTo(0, 2); got != 0 {
		t.Errorf("nothing scaled to %d", got)
	}

	// The limits go with it, since a minimum in points is what was meant.
	got := Limits{MinWidth: 640, MinHeight: 480, MaxWidth: 1280}.Scaled(1.5)
	if got.MinWidth != 960 || got.MinHeight != 720 || got.MaxWidth != 1920 {
		t.Errorf("the limits scaled to %+v", got)
	}
	if got.MaxHeight != 0 {
		t.Errorf("a bound that said nothing became %d", got.MaxHeight)
	}
}

func TestOpenWithRejectsASizeThatIsNotOne(t *testing.T) {
	if _, err := OpenWith(Options{Title: "x", Width: 0, Height: 100}); err == nil {
		t.Error("a window of no width opened")
	}
	if _, err := OpenWith(Options{Title: "x", Width: 100, Height: -1}); err == nil {
		t.Error("a window of negative height opened")
	}
}