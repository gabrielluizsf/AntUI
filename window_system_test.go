package antui

import "testing"

// TestSystemDarkIsReadEveryFrame is the answer the frame takes from the
// platform: it follows what the backend says each Begin, and the override —
// the stand-in for a platform that cannot ask — wins over it.
func TestSystemDarkIsReadEveryFrame(t *testing.T) {
	win, stub := newTestWindow(t, 64, 64)

	if dark, ok := win.SystemDark(); dark || ok {
		t.Errorf("a backend that has not said: SystemDark = %v, %v; want false, false", dark, ok)
	}

	stub.themeDark, stub.themeKnown = true, true
	win.Begin()
	if dark, ok := win.SystemDark(); !dark || !ok {
		t.Errorf("after a frame: SystemDark = %v, %v; want true, true — the frame asks the platform", dark, ok)
	}

	// An override is what a test and a platform with no sensor use, and it
	// holds however often the frame asks again.
	win.SetSystemDark(false)
	win.Begin()
	if dark, ok := win.SystemDark(); dark || !ok {
		t.Errorf("an override standing in for the platform: SystemDark = %v, %v; want false, true", dark, ok)
	}
	stub.themeDark, stub.themeKnown = true, true
	win.Begin()
	if dark, ok := win.SystemDark(); dark || !ok {
		t.Errorf("after another frame: SystemDark = %v, %v; want false, true — the override holds", dark, ok)
	}
}

// TestDisplayScaleIsTheRawAnswer is the difference between this and
// ContentScale: a display that never said has no resolution to report, while
// arithmetic on a scale still has a one to divide by.
func TestDisplayScaleIsTheRawAnswer(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)

	if scale, ok := win.DisplayScale(); scale != 0 || ok {
		t.Errorf("a display that has not said: DisplayScale = %v, %v; want 0, false", scale, ok)
	}
	if win.ContentScale() != 1 {
		t.Errorf("ContentScale = %v, want 1 — arithmetic never divides by nothing", win.ContentScale())
	}

	win.SetDisplayScale(1.5)
	if scale, ok := win.DisplayScale(); !ok || scale != 1.5 {
		t.Errorf("DisplayScale = %v, %v; want 1.5, true", scale, ok)
	}
	if win.ContentScale() != 1.5 {
		t.Errorf("ContentScale = %v, want 1.5 — the same number a query reads", win.ContentScale())
	}

	// Back to unanswered, which is what an offscreen window has all along.
	win.SetDisplayScale(0)
	if scale, ok := win.DisplayScale(); scale != 0 || ok {
		t.Errorf("after putting the answer back: DisplayScale = %v, %v; want 0, false", scale, ok)
	}
}

// TestDisplayScaleIsReadEveryFrame is the display's answer on the frame
// loop: what the display says lands on the frame after it, an override set
// by SetDisplayScale stands over it, and a display that stops answering does
// not take back what it said.
func TestDisplayScaleIsReadEveryFrame(t *testing.T) {
	win, stub := newTestWindow(t, 64, 64)

	if scale, ok := win.DisplayScale(); scale != 0 || ok {
		t.Errorf("a display that has not said: DisplayScale = %v, %v; want 0, false", scale, ok)
	}

	stub.scale = 1.5
	win.Begin()
	if scale, ok := win.DisplayScale(); !ok || scale != 1.5 {
		t.Errorf("after a frame: DisplayScale = %v, %v; want 1.5, true — the frame asks the display", scale, ok)
	}

	stub.scale = 2
	win.Begin()
	if scale, ok := win.DisplayScale(); !ok || scale != 2 {
		t.Errorf("after the display changed: DisplayScale = %v, %v; want 2, true", scale, ok)
	}

	// A display that stops saying leaves the number it gave: the scale
	// belongs to the screen, and a read with no answer is not the screen
	// taking it back.
	stub.scale = 0
	win.Begin()
	if scale, ok := win.DisplayScale(); !ok || scale != 2 {
		t.Errorf("after a display that said nothing: DisplayScale = %v, %v; want 2, true", scale, ok)
	}

	// An override is what a test and a platform with no sensor use, and it
	// holds however often the frame asks again.
	win.SetDisplayScale(1.25)
	stub.scale = 3
	win.Begin()
	if scale, ok := win.DisplayScale(); !ok || scale != 1.25 {
		t.Errorf("an override standing in for the display: DisplayScale = %v, %v; want 1.25, true", scale, ok)
	}
	win.Begin()
	if scale, ok := win.DisplayScale(); !ok || scale != 1.25 {
		t.Errorf("after another frame: DisplayScale = %v, %v; want 1.25, true — the override holds", scale, ok)
	}
}

// TestSystemDarkOnAWindowThatIsNotThere is the nil and closed forms, which a
// query can reach without an error to show for it.
func TestSystemDarkOnAWindowThatIsNotThere(t *testing.T) {
	var none *Window
	if dark, ok := none.SystemDark(); dark || ok {
		t.Errorf("SystemDark = %v, %v on no window; want false, false", dark, ok)
	}
	if scale, ok := none.DisplayScale(); scale != 0 || ok {
		t.Errorf("DisplayScale = %v, %v on no window; want 0, false", scale, ok)
	}
	none.SetSystemDark(true) // must not panic
	none.SetDisplayScale(2)  // must not panic
}
