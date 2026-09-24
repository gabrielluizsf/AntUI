package css

import (
	"math"
	"testing"
)

// TestTransitionParse pins a full shorthand down to its parts: property,
// duration, timing and delay, all the orderings CSS lets an author write.
func TestTransitionParse(t *testing.T) {
	tr, ok := parseTransition("background-color 0.5s cubic-bezier(0.25, 0.1, 0.25, 1) 0.1s")
	if !ok {
		t.Fatal("full shorthand should parse")
	}
	if tr.Prop != "background-color" {
		t.Errorf("prop = %q, want background-color", tr.Prop)
	}
	if tr.Duration.MS() != 500 {
		t.Errorf("duration = %vms, want 500", tr.Duration.MS())
	}
	if tr.Delay.MS() != 100 {
		t.Errorf("delay = %vms, want 100", tr.Delay.MS())
	}
	if tr.Timing.x1 != 0.25 || tr.Timing.y2 != 1 {
		t.Errorf("timing = %+v, want the written curve", tr.Timing)
	}

	reverse, ok := parseTransition("200ms opacity")
	if !ok {
		t.Fatal("time-first shorthand should parse")
	}
	if reverse.Prop != "opacity" || reverse.Duration.MS() != 200 {
		t.Errorf("time-first = %+v", reverse)
	}

	implicit, ok := parseTransition("1s")
	if !ok {
		t.Fatal("duration-only shorthand should parse")
	}
	if implicit.Prop != "all" || implicit.Duration.MS() != 1000 {
		t.Errorf("implicit all = %+v", implicit)
	}

	none, ok := parseTransition("none")
	if !ok {
		t.Fatal("none should parse")
	}
	if none.Prop != "none" {
		t.Errorf("none prop = %q", none.Prop)
	}
}

// TestTransitionParseInvalid drops a transition with no parts at all, and
// the whole shorthand when any entry is malformed. Two times are legal — the
// second is the delay.
func TestTransitionParseInvalid(t *testing.T) {
	if _, ok := parseTransition(""); ok {
		t.Error("empty transition should fail")
	}
	if list, ok := parseTransitions("1s, , 2s"); ok {
		t.Errorf("empty entry should fail the list, got %v", list)
	}
	durDelayed, ok := parseTransition("1s 2s")
	if !ok || durDelayed.Duration.MS() != 1000 || durDelayed.Delay.MS() != 2000 {
		t.Errorf("two times read as duration then delay, got %+v, %v", durDelayed, ok)
	}
}

// TestTransitionAllGaps a missing duration or timing means the transition
// still runs — instantly with the default easing — and a trailing value in
// any longhand spreads across the remaining entries.
func TestTransitionAllGaps(t *testing.T) {
	tr, ok := parseTransition("all 1s 0.3s")
	if !ok {
		t.Fatal("gap shorthand should parse")
	}
	if tr.Duration.MS() != 1000 || tr.Delay.MS() != 300 {
		t.Errorf("duration/delay = %v/%v", tr.Duration.MS(), tr.Delay.MS())
	}
	if tr.Timing != (Timing{}) {
		t.Error("default timing should be the zero Timing, which reads as linear")
	}
}

// TestTransitionLonghandsEachList the four longhands combine into a
// transition the way the shorthand would, repeating values across the list.
func TestTransitionLonghandsEachList(t *testing.T) {
	sh := mustParse(t, `
		button {
			transition-property: opacity, background-color;
			transition-duration: 1s;
			transition-timing-function: linear;
			transition-delay: 0.5s;
		}
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Transitions) != 2 {
		t.Fatalf("transitions = %d, want 2", len(st.Transitions))
	}
	if st.Transitions[0].Prop != "opacity" || st.Transitions[0].Duration.MS() != 1000 ||
		st.Transitions[0].Delay.MS() != 500 {
		t.Errorf("first transition = %+v", st.Transitions[0])
	}
	if st.Transitions[1].Prop != "background-color" || st.Transitions[1].Duration.MS() != 1000 {
		t.Errorf("second transition = %+v", st.Transitions[1])
	}
}

// TestTransitionNoneClears transition: none, however mixed with longhands,
// turns every transition off — the property that resets the cascade wins.
func TestTransitionNoneClears(t *testing.T) {
	sh := mustParse(t, `
		button {
			transition: none;
			transition-duration: 1s;
		}
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Transitions) != 0 {
		t.Errorf("transition: none should clear transitions, got %d", len(st.Transitions))
	}
}

// TestTransitionDurationsOnlyAll when only durations are written, the
// property list is empty and everything that follows animates ("all").
func TestTransitionDurationsOnlyAll(t *testing.T) {
	sh := mustParse(t, `button { transition-duration: 0.4s; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Transitions) != 1 || st.Transitions[0].Prop != "all" ||
		st.Transitions[0].Duration.MS() != 400 {
		t.Errorf("duration-only transitions = %+v", st.Transitions)
	}
}

// TestTransitionSpan the span of a transition is its delay plus duration,
// the number of milliseconds the drawer must run before it can settle.
func TestTransitionSpan(t *testing.T) {
	sh := mustParse(t, `
		button {
			transition: opacity 0.2s 0.1s, background-color 1s;
		}
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if got := TransitionSpan(st); got != 1000 {
		t.Errorf("TransitionSpan = %v, want 1000", got)
	}
	plain := mustParse(t, `button { color: red; }`).Style("button", nil, StateNone, 800)
	if TransitionSpan(plain) != 0 {
		t.Error("a style with no transitions spans zero")
	}
}

// TestTransitionShorthandStyle a transition written in the shorthand lands
// in the computed style the painter reads, with the easing applied.
func TestTransitionShorthandStyle(t *testing.T) {
	sh := mustParse(t, `
		button { transition: opacity 0.5s ease-in-out 0.2s; }
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Transitions) != 1 {
		t.Fatalf("transitions = %d, want 1", len(st.Transitions))
	}
	tr := st.Transitions[0]
	if tr.Prop != "opacity" || tr.Duration.MS() != 500 || tr.Delay.MS() != 200 {
		t.Errorf("transition = %+v", tr)
	}
	if got := tr.Timing.Ease(0.5); math.Abs(got-0.5) > 1e-6 {
		t.Errorf("ease-in-out at 0.5 = %v, want 0.5", got)
	}
}

// TestTransitionZeroDurationComposes a zero-duration transition is still a
// transition: it exists in the style, but the drawer snaps it.
func TestTransitionZeroDurationComposes(t *testing.T) {
	sh := mustParse(t, `button { transition: opacity 0s; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.Transitions) != 1 {
		t.Fatalf("transitions = %d, want 1", len(st.Transitions))
	}
	if st.Transitions[0].Duration.MS() != 0 {
		t.Errorf("duration = %vms, want 0", st.Transitions[0].Duration.MS())
	}
}
