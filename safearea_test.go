package antui

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// SafeArea keeps a program's pressable things out from under the status bar,
// the navigation bar and the notch. On a desktop it is the whole canvas; the
// interesting parts are the clamping, because the platform reports this in
// its own coordinates and a stale one after a resize would otherwise send
// drawing off the end of the canvas.

func TestSafeAreaIsEverythingByDefault(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	got := win.SafeArea()
	want := canvas.Area{X: 0, Y: 0, Width: 320, Height: 200}
	if got != want {
		t.Errorf("default safe area is %+v, want %+v", got, want)
	}
}

// A safe area hangs off the edge of the canvas? The window owns every pixel;
// the canvas is the edge, and the area is clamped to it.
func TestSafeAreaClampsToTheCanvas(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.SetSafeArea(canvas.Area{X: 10, Y: 20, Width: 400, Height: 300})
	got := win.SafeArea()
	want := canvas.Area{X: 10, Y: 20, Width: 310, Height: 180}
	if got != want {
		t.Errorf("clamped safe area is %+v, want %+v", got, want)
	}
}

// A negative inset is a status bar that has moved beyond the corner; the
// area is pulled back to the edge it is crossing, not widened to cover it.
func TestSafeAreaPullsBackNegativeInsets(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.SetSafeArea(canvas.Area{X: -5, Y: -10, Width: 320, Height: 200})
	got := win.SafeArea()
	if got.X != 0 || got.Y != 0 {
		t.Errorf("a negative inset should be pulled back to the edge, got %+v", got)
	}
	if got.Width > 320 {
		t.Errorf("a negative left inset left the area %d wide", got.Width)
	}
}

// A safe area that does not even touch the canvas after clamping is the
// platform talking about a stale size; the window falls back to "the whole
// canvas" rather than drawing nothing somewhere.
func TestSafeAreaOffTheCanvasFallsBackToEverything(t *testing.T) {
	win, _ := newTestWindow(t, 320, 200)
	win.SetSafeArea(canvas.Area{X: 400, Y: 190, Width: 50, Height: 50})
	got := win.SafeArea()
	want := canvas.Area{X: 0, Y: 0, Width: 320, Height: 200}
	if got != want {
		t.Errorf("a safe area off the canvas is %+v, want %+v", got, want)
	}
}