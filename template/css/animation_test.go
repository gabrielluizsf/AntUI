package css

import (
	"math"
	"testing"
)

// TestAnimationShorthandComposes a full shorthand reads into one animation
// entry with every knob set, and lands in the computed style.
func TestAnimationShorthandComposes(t *testing.T) {
	list, ok := parseAnimations("spin 2s cubic-bezier(0.25, 0.1, 0.25, 1) 0.5s infinite alternate-reverse both")
	if !ok || len(list) != 1 {
		t.Fatalf("parseAnimations = %v, %v", list, ok)
	}
	a := list[0]
	if a.Name != "spin" || a.Duration.MS() != 2000 || a.Delay.MS() != 500 {
		t.Errorf("name/duration/delay = %q/%vms/%vms", a.Name, a.Duration.MS(), a.Delay.MS())
	}
	if !a.Infinite {
		t.Error("infinite was not read")
	}
	if a.Direction != AnimAlternateReverse {
		t.Errorf("direction = %v, want alternate-reverse", a.Direction)
	}
	if a.Fill != FillBoth {
		t.Errorf("fill = %v, want both", a.Fill)
	}

	sh := mustParse(t, `button { animation: spin 2s linear 0.5s infinite alternate both; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Animations) != 1 {
		t.Fatalf("animations = %d, want 1", len(st.Animations))
	}
	composed := st.Animations[0]
	if composed.Name != "spin" || !composed.Infinite || composed.Direction != AnimAlternate {
		t.Errorf("composed animation = %+v", composed)
	}
}

// TestAnimationNamesInAnyOrder an animation's parts appear in any order, and
// both a keyword and a quoted string make a name.
func TestAnimationNamesInAnyOrder(t *testing.T) {
	a, ok := parseAnimation("2s spin reverse")
	if !ok {
		t.Fatal("parts out of order should parse")
	}
	if a.Name != "spin" || a.Duration.MS() != 2000 || a.Direction != AnimReverse {
		t.Errorf("out-of-order animation = %+v", a)
	}
	quoted, ok := parseAnimation(`"fade-in" 1s`)
	if !ok || quoted.Name != "fade-in" {
		t.Errorf("quoted name = %+v, %v", quoted, ok)
	}
}

// TestAnimationIterationFinite the finite iteration count defaults to one
// and a written number is kept, not mistaken for a duration.
func TestAnimationIterationFinite(t *testing.T) {
	a, ok := parseAnimation("bounce 3")
	if !ok {
		t.Fatal("bounce 3 should parse")
	}
	if a.Iterations != 3 || a.Infinite {
		t.Errorf("iterations = %v infinite=%v, want 3", a.Iterations, a.Infinite)
	}
	sh := mustParse(t, `button { animation: bounce 2s 3; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if st.Animations[0].Iterations != 3 {
		t.Errorf("composed iterations = %v, want 3", st.Animations[0].Iterations)
	}
}

// TestAnimationDefaultsGaps a naming-only animation fills its longhands with
// the CSS initials: zero duration, the default timing, no delay, one
// iteration, normal direction, none fill.
func TestAnimationDefaultsGaps(t *testing.T) {
	a, ok := parseAnimation("float")
	if !ok {
		t.Fatal("name-only animation should parse")
	}
	if a.Duration.MS() != 0 || a.Iterations != 1 || a.Direction != AnimNormal || a.Fill != FillNone {
		t.Errorf("defaults = %+v", a)
	}
}

// TestAnimationNoneInvalid an animation with no name is not an animation: the
// declaration is dropped the way a browser drops an unparsable one.
func TestAnimationNoneInvalid(t *testing.T) {
	if _, ok := parseAnimation("2s"); ok {
		t.Error("duration without a name should fail")
	}
	if _, ok := parseAnimations("infinite"); ok {
		t.Error("keywords without a name should fail")
	}
}

// TestAnimationMultipleGapFill a shorter name list spreads the longer
// longhands across the animations, exactly as CSS repeats the short list.
func TestAnimationMultipleGapFill(t *testing.T) {
	sh := mustParse(t, `
		button {
			animation-name: a, b, c;
			animation-duration: 1s, 2s, 3s;
			animation-iteration-count: infinite;
			animation-direction: alternate;
		}
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Animations) != 3 {
		t.Fatalf("animations = %d, want 3", len(st.Animations))
	}
	if !st.Animations[2].Infinite || st.Animations[1].Direction != AnimAlternate ||
		st.Animations[0].Duration.MS() != 1000 || st.Animations[2].Duration.MS() != 3000 {
		t.Errorf("gap-filled animations = %+v", st.Animations)
	}
}

// TestAnimationInfiniteFlowsIntoList an infinite keyword lands as positive
// infinity in the parallel iteration list, and a finite count as the number.
func TestAnimationInfiniteFlowsIntoList(t *testing.T) {
	names, durs, tims, dels, iters, dirs, fills := animationParts(
		[]Animation{{Name: "a", Infinite: true}, {Name: "b", Iterations: 2}},
	)
	if len(names) != 2 || !math.IsInf(iters[0], 1) || iters[1] != 2 {
		t.Errorf("parts = %v %v %v %v %v %v %v", names, durs, tims, dels, iters, dirs, fills)
	}
}

// TestAnimationNoKeyframesStillComposes an animation whose @keyframes block
// was never written still sits in the style; the drawer finds no block and
// skips it, the way a browser skips a missing animation.
func TestAnimationNoKeyframesStillComposes(t *testing.T) {
	sh := mustParse(t, `button { animation: ghost 1s; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Animations) != 1 || st.Animations[0].Name != "ghost" {
		t.Errorf("animations = %+v", st.Animations)
	}
}
