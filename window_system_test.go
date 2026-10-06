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
	none.adoptSystemTheme()  // must not panic
}

// TestWindowStartsWithTheSystemTheme is the half of the system's answer that
// is a theme rather than a media query: what the system paints in is the
// theme a window opens with, so a program that never chooses one is not a
// light window on a dark desktop. Only dark is written over what the window
// already had — it opens light — and a system that said nothing says
// nothing here too.
func TestWindowStartsWithTheSystemTheme(t *testing.T) {
	for _, want := range []struct {
		name         string
		dark, known  bool
		isDarkSystem bool
	}{
		{name: "a dark system", dark: true, known: true, isDarkSystem: true},
		{name: "a light system", dark: false, known: true},
		{name: "a system that said nothing", dark: false, known: false},
		{name: "a system that said dark but is not known to", dark: true, known: false},
	} {
		win, stub := newTestWindow(t, 64, 64)
		stub.themeDark, stub.themeKnown = want.dark, want.known
		win.adoptSystemTheme()

		switch got := *win.Theme(); {
		case want.isDarkSystem && got != DarkTheme():
			t.Errorf("%s: theme = %#v, want the dark one the system paints in", want.name, got)
		case !want.isDarkSystem && got != LightTheme():
			t.Errorf("%s: theme = %#v, want the light one the window opened with", want.name, got)
		}
	}

	// From the start the theme is the program's: SetTheme is the last word
	// over what the system began it with, and a frame's worth of system
	// answers goes on without restyling the window underneath it.
	win, stub := newTestWindow(t, 64, 64)
	stub.themeDark, stub.themeKnown = true, true
	win.adoptSystemTheme()
	win.SetTheme(LightTheme())
	win.Begin()
	if *win.Theme() != LightTheme() {
		t.Errorf("theme after a frame = %#v, want the light one the program chose", *win.Theme())
	}

	// And a window with no system behind it — offscreen, or closed — has
	// nothing to take and does not fail trying.
	off, _, err := Offscreen(8, 8)
	if err != nil {
		t.Fatalf("Offscreen: %v", err)
	}
	off.adoptSystemTheme()
	if *off.Theme() != LightTheme() {
		t.Errorf("offscreen theme = %#v, want the light one", *off.Theme())
	}
}
