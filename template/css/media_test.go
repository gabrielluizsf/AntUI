package css

import (
	"strings"
	"testing"
)

// mediaWidth reads one class's width against a whole window, the way the
// drawing asks for it, and reports the rule as absent when it did not apply.
func mediaWidth(sh *Sheet, vp Viewport, cls string) (int, bool) {
	st := sh.StyleViewport("div", []string{cls}, StateNone, vp)
	if !st.Has("width") {
		return 0, false
	}
	return st.Width.Px(vp.Width), true
}

// TestMediaReadsOrientation checks the window's shape against the query: a
// window at least as tall as it is wide is portrait, the other way round is
// landscape, and a condition with no value the spec gives still leaves the
// rule standing.
func TestMediaReadsOrientation(t *testing.T) {
	sh := parseOK(t, `
@media (orientation: portrait)  { .p { width: 10px; } }
@media (orientation: landscape) { .l { width: 20px; } }
@media not (orientation: portrait) { .n { width: 30px; } }
@media (orientation: sideways) { .bad { width: 40px; } }
`)
	tall := Viewport{Width: 400, Height: 800}
	wide := Viewport{Width: 800, Height: 400}
	square := Viewport{Width: 600, Height: 600}

	if v, ok := mediaWidth(sh, tall, "p"); !ok || v != 10 {
		t.Errorf("portrait at 400x800: %d, %v; want 10, true", v, ok)
	}
	if _, ok := mediaWidth(sh, wide, "p"); ok {
		t.Error("portrait should reject a window twice as wide as it is tall")
	}
	if v, ok := mediaWidth(sh, square, "p"); !ok || v != 10 {
		t.Errorf("portrait on a square window: %d, %v; want 10, true — the "+
			"height reaching the width counts", v, ok)
	}
	if v, ok := mediaWidth(sh, wide, "l"); !ok || v != 20 {
		t.Errorf("landscape at 800x400: %d, %v; want 20, true", v, ok)
	}
	if _, ok := mediaWidth(sh, tall, "l"); ok {
		t.Error("landscape should reject a tall window")
	}
	if _, ok := mediaWidth(sh, square, "l"); ok {
		t.Error("landscape should reject a square window")
	}
	if v, ok := mediaWidth(sh, wide, "n"); !ok || v != 30 {
		t.Errorf("not-portrait at 800x400: %d, %v; want 30, true", v, ok)
	}
	if _, ok := mediaWidth(sh, tall, "n"); ok {
		t.Error("not-portrait should reject the portrait window")
	}
	if v, ok := mediaWidth(sh, tall, "bad"); !ok || v != 40 {
		t.Errorf("an orientation the spec does not give: %d, %v; want 40, "+
			"true — a condition with no answer never drops the rule", v, ok)
	}
	if len(sh.Warn) != 1 || !strings.Contains(sh.Warn[0], "orientation") {
		t.Errorf("sh.Warn = %q, want one about the orientation", sh.Warn)
	}
}

// TestMediaReadsResolution reads the display's density, in whichever unit a
// stylesheet writes it, and leaves the answer alone when the display gave
// none.
func TestMediaReadsResolution(t *testing.T) {
	sh := parseOK(t, `
@media (min-resolution: 2dppx)  { .hi { width: 10px; } }
@media (resolution: 96dpi)      { .one { width: 20px; } }
@media (max-resolution: 1dppx)  { .lo { width: 30px; } }
@media (min-resolution: 40dpcm) { .cm { width: 40px; } }
@media (min-resolution: 3 dots) { .bad { width: 50px; } }
`)
	retina := Viewport{Width: 800, Height: 600, Scale: 2}
	plain := Viewport{Width: 800, Height: 600, Scale: 1}
	unsaid := Viewport{Width: 800, Height: 600}

	if v, ok := mediaWidth(sh, retina, "hi"); !ok || v != 10 {
		t.Errorf("2dppx on a display drawing two pixels a point: %d, %v; want 10, true", v, ok)
	}
	if _, ok := mediaWidth(sh, plain, "hi"); ok {
		t.Error("2dppx should reject a display drawing one pixel a point")
	}
	if v, ok := mediaWidth(sh, plain, "one"); !ok || v != 20 {
		t.Errorf("resolution: 96dpi at scale 1: %d, %v; want 20, true", v, ok)
	}
	if _, ok := mediaWidth(sh, retina, "one"); ok {
		t.Error("resolution: 96dpi should reject a display at twice that")
	}
	if v, ok := mediaWidth(sh, plain, "lo"); !ok || v != 30 {
		t.Errorf("max-resolution 1dppx at scale 1: %d, %v; want 30, true", v, ok)
	}
	if _, ok := mediaWidth(sh, retina, "lo"); ok {
		t.Error("max-resolution 1dppx should reject a display at 2dppx")
	}
	if _, ok := mediaWidth(sh, plain, "cm"); ok {
		t.Error("40dpcm is 101.6dpi, which a 96dpi display does not reach")
	}
	if v, ok := mediaWidth(sh, unsaid, "cm"); !ok || v != 40 {
		t.Errorf("a display that said nothing: %d, %v; want 40, true — the "+
			"query is neither kept out nor forced", v, ok)
	}
	if v, ok := mediaWidth(sh, unsaid, "hi"); !ok || v != 10 {
		t.Errorf("min-resolution on a display that said nothing: %d, %v; want 10, true", v, ok)
	}
	if len(sh.Warn) != 1 || !strings.Contains(sh.Warn[0], "resolution") {
		t.Errorf("sh.Warn = %q, want one about the resolution", sh.Warn)
	}
}

// TestMediaReadsColorScheme puts the system's own scheme under the query,
// including the system that has no answer to give.
func TestMediaReadsColorScheme(t *testing.T) {
	sh := parseOK(t, `
@media (prefers-color-scheme: dark)  { .d { width: 10px; } }
@media (prefers-color-scheme: light) { .lt { width: 20px; } }
@media not (prefers-color-scheme: dark) { .n { width: 30px; } }
@media (prefers-color-scheme: sepia) { .bad { width: 40px; } }
`)
	dark := Viewport{Width: 400, Height: 400, Scheme: SchemeDark}
	light := Viewport{Width: 400, Height: 400, Scheme: SchemeLight}
	unsaid := Viewport{Width: 400, Height: 400}

	if v, ok := mediaWidth(sh, dark, "d"); !ok || v != 10 {
		t.Errorf("dark on a dark system: %d, %v; want 10, true", v, ok)
	}
	if _, ok := mediaWidth(sh, light, "d"); ok {
		t.Error("prefers-color-scheme: dark should reject a light system")
	}
	if v, ok := mediaWidth(sh, light, "lt"); !ok || v != 20 {
		t.Errorf("light on a light system: %d, %v; want 20, true", v, ok)
	}
	if v, ok := mediaWidth(sh, light, "n"); !ok || v != 30 {
		t.Errorf("not dark on a light system: %d, %v; want 30, true", v, ok)
	}
	if _, ok := mediaWidth(sh, dark, "n"); ok {
		t.Error("not dark should reject a dark system")
	}
	if v, ok := mediaWidth(sh, unsaid, "d"); !ok || v != 10 {
		t.Errorf("dark on a system that said nothing: %d, %v; want 10, true — "+
			"the rule is never dropped for want of an answer", v, ok)
	}
	if v, ok := mediaWidth(sh, unsaid, "n"); !ok || v != 30 {
		t.Errorf("not dark on a system that said nothing: %d, %v; want 30, true", v, ok)
	}
	if v, ok := mediaWidth(sh, dark, "bad"); !ok || v != 40 {
		t.Errorf("a scheme the spec does not give: %d, %v; want 40, true", v, ok)
	}
	if len(sh.Warn) != 1 || !strings.Contains(sh.Warn[0], "color-scheme") {
		t.Errorf("sh.Warn = %q, want one about the color-scheme", sh.Warn)
	}
}

// TestMediaNestedQueriesDisagree walks the intersection of two nested @media:
// queries that agree carry over, and queries that cannot both hold leave the
// rule belonging to no window at all.
func TestMediaNestedQueriesDisagree(t *testing.T) {
	sh := parseOK(t, `
@media (orientation: portrait) {
  @media (orientation: portrait) { .agree { width: 10px; } }
  @media (orientation: landscape) { .shape { width: 20px; } }
}
@media (prefers-color-scheme: dark) {
  @media (prefers-color-scheme: dark) { .scheme { width: 30px; } }
  @media (prefers-color-scheme: light) { .colour { width: 40px; } }
}
@media (orientation: landscape) {
  @media (orientation: landscape) and (min-resolution: 2dppx) { .both { width: 50px; } }
}
`)
	tall := Viewport{Width: 400, Height: 800, Scheme: SchemeDark, Scale: 1}
	wide := Viewport{Width: 800, Height: 400, Scheme: SchemeDark, Scale: 2}

	if v, ok := mediaWidth(sh, tall, "agree"); !ok || v != 10 {
		t.Errorf("portrait inside portrait: %d, %v; want 10, true", v, ok)
	}
	if _, ok := mediaWidth(sh, tall, "shape"); ok {
		t.Error("a landscape window inside a portrait one matches nothing")
	}
	if v, ok := mediaWidth(sh, tall, "scheme"); !ok || v != 30 {
		t.Errorf("dark inside dark: %d, %v; want 30, true", v, ok)
	}
	if _, ok := mediaWidth(sh, tall, "colour"); ok {
		t.Error("a light query inside a dark one matches nothing")
	}
	if v, ok := mediaWidth(sh, wide, "both"); !ok || v != 50 {
		t.Errorf("landscape and 2dppx where both hold: %d, %v; want 50, true", v, ok)
	}
	if _, ok := mediaWidth(sh, tall, "both"); ok {
		t.Error("the landscape half of the nested pair should keep it out")
	}
	if _, ok := mediaWidth(sh, Viewport{Width: 800, Height: 400, Scheme: SchemeDark, Scale: 1}, "both"); ok {
		t.Error("the resolution half of the nested pair should keep it out")
	}
}

// TestMediaNestedQueriesNarrow is the overlap two nested @media leave:
// whichever edge both windows declare — pixels or dots per inch, a minimum
// or a maximum — meets at the tighter of the two, and an edge only one of
// them declares passes through untouched.
func TestMediaNestedQueriesNarrow(t *testing.T) {
	sh := parseOK(t, `
@media (min-width: 400px) {
  @media (min-width: 700px) { .hi { width: 10px; } }
}
@media (max-width: 700px) {
  @media (max-width: 400px) { .lo { width: 20px; } }
}
@media (min-resolution: 2dppx) {
  @media (min-resolution: 3dppx) { .retina { width: 30px; } }
}
@media (max-resolution: 2dppx) {
  @media (max-resolution: 1.5dppx) { .small { width: 40px; } }
}
@media (min-width: 400px) {
  @media (min-height: 300px) { .cross { width: 50px; } }
}
@media (min-width: 700px) {
  @media (min-width: 400px) { .looser { width: 60px; } }
}
`)
	at := func(vp Viewport, cls string) (int, bool) { return mediaWidth(sh, vp, cls) }

	// Two minima meet at the higher one, two maxima at the lower one.
	if v, ok := at(Viewport{Width: 800, Height: 400}, "hi"); !ok || v != 10 {
		t.Errorf("800 wide inside a 400 and 700: %d, %v; want 10, true", v, ok)
	}
	if _, ok := at(Viewport{Width: 500, Height: 400}, "hi"); ok {
		t.Error("the inner 700 should keep a 500 wide window out")
	}
	if v, ok := at(Viewport{Width: 400, Height: 400}, "lo"); !ok || v != 20 {
		t.Errorf("400 wide inside a 700 and 400: %d, %v; want 20, true", v, ok)
	}
	if _, ok := at(Viewport{Width: 500, Height: 400}, "lo"); ok {
		t.Error("the inner 400 should keep a 500 wide window out")
	}

	// The same two edges in dots per inch: 3dppx is 288dpi, 1.5dppx is 144.
	if v, ok := at(Viewport{Width: 400, Height: 400, Scale: 3}, "retina"); !ok || v != 30 {
		t.Errorf("288dpi against a 192 and 288 floor: %d, %v; want 30, true", v, ok)
	}
	if _, ok := at(Viewport{Width: 400, Height: 400, Scale: 2}, "retina"); ok {
		t.Error("the inner 3dppx should keep a 2dppx display out")
	}
	if v, ok := at(Viewport{Width: 400, Height: 400, Scale: 1}, "small"); !ok || v != 40 {
		t.Errorf("96dpi against a 192 and 144 ceiling: %d, %v; want 40, true", v, ok)
	}
	if _, ok := at(Viewport{Width: 400, Height: 400, Scale: 2}, "small"); ok {
		t.Error("the inner 1.5dppx should keep a 2dppx display out")
	}

	// An edge only one of the two windows declares is not narrowed at all.
	if v, ok := at(Viewport{Width: 800, Height: 400}, "cross"); !ok || v != 50 {
		t.Errorf("800 by 400 inside a width and a height: %d, %v; want 50, true", v, ok)
	}
	if _, ok := at(Viewport{Width: 800, Height: 200}, "cross"); ok {
		t.Error("the inner 300 should keep a 200 tall window out")
	}
	if _, ok := at(Viewport{Width: 300, Height: 400}, "cross"); ok {
		t.Error("the outer 400 should keep a 300 wide window out")
	}

	// A nested window that asks for less than the one around it does not
	// loosen what already holds: the outer 700 still decides.
	if v, ok := at(Viewport{Width: 800, Height: 400}, "looser"); !ok || v != 60 {
		t.Errorf("800 wide inside a 700 and 400: %d, %v; want 60, true", v, ok)
	}
	if _, ok := at(Viewport{Width: 500, Height: 400}, "looser"); ok {
		t.Error("the outer 700 should keep a 500 wide window out")
	}
}

// TestMediaResolutionRefusesWhatItCannotRead is the parser meeting a
// resolution it cannot make sense of: no digits behind the number, or no
// unit the spec gives one in, is a warning and a condition left out of the
// answer rather than a number half read.
func TestMediaResolutionRefusesWhatItCannotRead(t *testing.T) {
	sh := parseOK(t, `
@media (max-resolution: ..dpi) { .unitless { width: 10px; } }
@media (min-resolution: 300) { .bare { width: 20px; } }
`)
	if len(sh.Warn) != 2 {
		t.Fatalf("sh.Warn = %q, want one for each value", sh.Warn)
	}
	for _, w := range sh.Warn {
		if !strings.Contains(w, "resolution") {
			t.Errorf("sh.Warn holds %q, want it to name the resolution", w)
		}
	}
	// Neither condition has an edge to hold, so both rules stand.
	if v, ok := mediaWidth(sh, Viewport{Width: 400, Height: 400, Scale: 1}, "unitless"); !ok || v != 10 {
		t.Errorf("a resolution with no digits: %d, %v; want 10, true", v, ok)
	}
	if v, ok := mediaWidth(sh, Viewport{Width: 400, Height: 400}, "bare"); !ok || v != 20 {
		t.Errorf("a resolution with no unit: %d, %v; want 20, true", v, ok)
	}
}
