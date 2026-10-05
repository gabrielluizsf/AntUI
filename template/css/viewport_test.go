package css

import "testing"

// TestMediaMeasuresWindowHeight walks a sheet's @media answers across a
// window's two edges: a query that only names a height has to see the height,
// and one that names both has to see both.
func TestMediaMeasuresWindowHeight(t *testing.T) {
	sh := parseOK(t, `
@media (min-height: 600px) { .tall  { width: 10px; } }
@media (max-height: 500px) { .short { width: 20px; } }
@media (min-width: 700px) and (min-height: 600px) { .both { width: 30px; } }
@media (min-width: 700px) { .wide { width: 40px; } }
`)
	widthAt := func(vp Viewport, cls string) (int, bool) {
		st := sh.StyleViewport("div", []string{cls}, StateNone, vp)
		if !st.Has("width") {
			return 0, false
		}
		return st.Width.Px(vp.Width), true
	}

	if v, ok := widthAt(Viewport{Width: 800, Height: 800}, "tall"); !ok || v != 10 {
		t.Errorf("min-height at 800x800: %d, %v; want 10, true", v, ok)
	}
	if _, ok := widthAt(Viewport{Width: 800, Height: 400}, "tall"); ok {
		t.Error("min-height: 600px should reject a 400 tall window")
	}
	if v, ok := widthAt(Viewport{Width: 800, Height: 400}, "short"); !ok || v != 20 {
		t.Errorf("max-height at 800x400: %d, %v; want 20, true", v, ok)
	}
	if _, ok := widthAt(Viewport{Width: 800, Height: 800}, "short"); ok {
		t.Error("max-height: 500px should reject an 800 tall window")
	}

	if v, ok := widthAt(Viewport{Width: 800, Height: 800}, "both"); !ok || v != 30 {
		t.Errorf("and-window at 800x800: %d, %v; want 30, true", v, ok)
	}
	if _, ok := widthAt(Viewport{Width: 800, Height: 400}, "both"); ok {
		t.Error("and-window should reject a window too short for its height half")
	}
	if _, ok := widthAt(Viewport{Width: 500, Height: 800}, "both"); ok {
		t.Error("and-window should reject a window too narrow for its width half")
	}
	if v, ok := widthAt(Viewport{Width: 800, Height: 400}, "wide"); !ok || v != 40 {
		t.Errorf("width query at 800x400: %d, %v; want 40, true; a height "+
			"the query never names must not weigh in", v, ok)
	}
}

// TestStyleResolvesViewportUnitsAgainstTheWindow pins the measurement context
// the cascade hands the length parser: a square stand-in window answers a
// viewport unit off its own width, and only a real window answers off its
// height.
func TestStyleResolvesViewportUnitsAgainstTheWindow(t *testing.T) {
	sh := parseOK(t, "button { width: calc(10vh); }")

	square := sh.Style("button", nil, StateNone, 400)
	if got := square.Width.Px(400); got != 40 {
		t.Errorf("width under the width-only call: %d, want 40 — a window as "+
			"wide as it is tall", got)
	}

	real := sh.StyleViewport("button", nil, StateNone, Viewport{Width: 400, Height: 800})
	if got := real.Width.Px(400); got != 80 {
		t.Errorf("width under an 800 tall window: %d, want 80 — 10vh of the "+
			"window it is drawn in", got)
	}
}

// TestGetStyleViewportCacheTracksWindow checks that the table's cached styles
// belong to the window they were computed for: a second window of the same
// width but a different height is a different answer, and nothing of the first
// one is left behind to be handed out.
func TestGetStyleViewportCacheTracksWindow(t *testing.T) {
	c := load(t, `
		button { background-color: #00FF00; }
		@media (max-height: 500px) { button { background-color: #FF0000; } }
	`)
	short := Viewport{Width: 400, Height: 400}
	tall := Viewport{Width: 400, Height: 900}

	st := c.GetStyleViewport(RoleButton, nil, StateNone, short)
	if st.Background.R() != 0xFF || st.Background.G() != 0x00 {
		t.Errorf("short window: background %v, want the red max-height rule", st.Background)
	}
	if len(c.specs) == 0 {
		t.Error("expected the answer to have been cached")
	}

	st = c.GetStyleViewport(RoleButton, nil, StateNone, tall)
	if st.Background.G() != 0xFF || st.Background.R() != 0x00 {
		t.Errorf("tall window: background %v, want the green rule — the red "+
			"answer cached for another window leaked", st.Background)
	}
	for k := range c.specs {
		if k.width != tall.Width || k.height != tall.Height {
			t.Errorf("cache holds a key for %dx%d, want only %dx%d",
				k.width, k.height, tall.Width, tall.Height)
		}
	}
}
