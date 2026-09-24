package template

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/template/css"
)

// TestCSSAnimationLoops an infinite animation turns back to its start when
// its duration passes: at a full cycle and a half the box is halfway from
// from to to again.
func TestCSSAnimationLoops(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { animation: fade 1s linear infinite; }
		@keyframes fade { from { opacity: 0; } to { opacity: 1; } }
	`)

	*clock = 0
	_, _, _, _, begin, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be visible")
	}
	if begin.Opacity != 0 {
		t.Errorf("opacity at 0s = %v, want 0", begin.Opacity)
	}

	*clock = 0.5
	half := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if math.Abs(half.Opacity-0.5) > 1e-9 {
		t.Errorf("opacity at 0.5s = %v, want 0.5", half.Opacity)
	}

	*clock = 0.9
	nine := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if math.Abs(nine.Opacity-0.9) > 1e-9 {
		t.Errorf("opacity at 0.9s = %v, want 0.9", nine.Opacity)
	}

	// 1.5s starts the second lap; progress is back to the middle.
	*clock = 1.5
	lapped := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if math.Abs(lapped.Opacity-0.5) > 1e-9 {
		t.Errorf("opacity at 1.5s = %v, want 0.5", lapped.Opacity)
	}
}

// TestCSSAnimationForwardFill a finished animation keeps its last frame when
// fill-mode says so, and waits through its delay before it starts.
func TestCSSAnimationForwardFill(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { animation: fade 1s linear 0.5s forwards; }
		@keyframes fade { from { opacity: 0; } to { opacity: 1; } }
	`)

	// Inside the delay nothing has played yet and there is no backwards
	// fill, so the base (unset) opacity stays.
	*clock = 0
	_, _, _, _, base, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be visible")
	}
	if base.Opacity != 0 {
		t.Errorf("opacity during the delay = %v, want 0", base.Opacity)
	}

	// A quarter of a second into the duration.
	*clock = 0.75
	playing := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if math.Abs(playing.Opacity-0.25) > 1e-9 {
		t.Errorf("opacity during playback = %v, want 0.25", playing.Opacity)
	}

	// Past the end the last frame sticks because of the forwards fill.
	*clock = 2
	pinned := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if pinned.Opacity != 1 {
		t.Errorf("opacity after the end = %v, want 1", pinned.Opacity)
	}
}

// TestCSSAnimationAlternate an alternate-direction animation travels back
// down on its second lap instead of popping to the start.
func TestCSSAnimationAlternate(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { animation: fade 1s linear alternate infinite; }
		@keyframes fade { from { opacity: 0; } to { opacity: 1; } }
	`)

	*clock = 0
	tpl.layout(css.RoleButton, "Go")
	*clock = 1.25
	down := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if math.Abs(down.Opacity-0.75) > 1e-9 {
		t.Errorf("opacity on the return lap = %v, want 0.75", down.Opacity)
	}
}

// TestCSSAnimationKeyframesGlide a background-colour animation glides the
// two frame colours instead of swapping, exactly like a transition between
// the same pair.
func TestCSSAnimationKeyframesGlide(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; animation: glow 1s linear infinite; }
		@keyframes glow { from { background-color: #111111; } to { background-color: #333333; } }
	`)

	*clock = 0.5
	_, _, _, _, mid, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be visible")
	}
	if got := mid.Background; got != 0xFF222222 {
		t.Errorf("background halfway = %v, want #222222", got)
	}
}
