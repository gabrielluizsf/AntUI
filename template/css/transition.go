package css

import (
	"strings"
)

// Transition is one property an element animates between its computed values
// when the cascade changes it: the transition shorthand and its longhands
// come together here. Prop is the canonical property name (such as "opacity"
// or "background-color") or "all"; Duration and Delay are the two times CSS
// writes in either order; Timing is the easing. A zero Duration means the
// property snaps, exactly as a style without a transition does.
type Transition struct {
	Prop     string
	Duration Time
	Timing   Timing
	Delay    Time
}

// parseTransitions reads the transition shorthand: a comma-separated list of
// transitions. "none" clears every entry, so a style returns no transitions
// at all. An empty or unrecognised entry invalidates the whole property, the
// way a browser drops a declaration it cannot parse.
func parseTransitions(raw string) ([]Transition, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if raw == "none" {
		return nil, true
	}
	var out []Transition
	for _, part := range splitFields(raw, ',') {
		tr, ok := parseTransition(strings.TrimSpace(part))
		if !ok {
			return nil, false
		}
		out = append(out, tr)
	}
	return out, true
}

// parseTransition reads one comma-free transition: a property name (which may
// be "all" or "none"), two times with the duration first, an easing, in any
// order. A missing property means every property, and a missing timing means
// the default silence of [Linear].
func parseTransition(part string) (Transition, bool) {
	part = strings.TrimSpace(part)
	if part == "" {
		return Transition{}, false
	}
	var tr Transition
	var sawDur, sawDelay bool
	for _, tk := range splitTokens(part) {
		if tim, ok := parseTiming(tk); ok {
			tr.Timing = tim
			continue
		}
		if tm, ok := ParseTime(tk); ok {
			if !sawDur {
				tr.Duration = tm
				sawDur = true
			} else if !sawDelay {
				tr.Delay = tm
				sawDelay = true
			} else {
				return Transition{}, false
			}
			continue
		}
		if tr.Prop == "" {
			if tk == "none" {
				tr.Prop = "none"
			} else if !strings.ContainsRune(tk, '(') {
				tr.Prop = strings.ToLower(tk)
			}
		}
	}
	if tr.Prop == "" {
		tr.Prop = "all"
	}
	return tr, true
}

// transitionParts splits a parsed shorthand into the four parallel lists the
// longhands also fill, so the cascade can merge them all in one finisher.
func transitionParts(ts []Transition) (props []string, durs []Time, tims []Timing, dels []Time) {
	for _, tr := range ts {
		props = append(props, tr.Prop)
		durs = append(durs, tr.Duration)
		tims = append(tims, tr.Timing)
		dels = append(dels, tr.Delay)
	}
	return props, durs, tims, dels
}

// finishTransitions merges the transition longhands (and the shorthand's own
// parallel lists) into the style's final Transitions slice. Like CSS, every
// list repeats cyclically to cover the whole row, so a single duration drives
// every named property; the property list sets the row's length, and a row of
// timings alone — durations without a property, which mean "all" — counts as
// one transition.
func finishTransitions(st *Style) {
	if st.TransitionNone {
		st.Transitions = nil
		return
	}
	n := len(st.TransitionProps)
	onlyTimings := len(st.TransitionDurs) == 0 &&
		len(st.TransitionTims) == 0 && len(st.TransitionDels) == 0
	if n == 0 {
		if onlyTimings {
			st.Transitions = nil
			return
		}
		n = 1 // durations, easings or delays without a property: "all"
	}
	for i := 0; i < n; i++ {
		tr := Transition{Prop: "all"}
		if i < len(st.TransitionProps) {
			tr.Prop = st.TransitionProps[i]
		}
		if tr.Prop == "none" {
			continue
		}
		if len(st.TransitionDurs) > 0 {
			tr.Duration = st.TransitionDurs[i%len(st.TransitionDurs)]
		}
		if len(st.TransitionTims) > 0 {
			tr.Timing = st.TransitionTims[i%len(st.TransitionTims)]
		}
		if len(st.TransitionDels) > 0 {
			tr.Delay = st.TransitionDels[i%len(st.TransitionDels)]
		}
		st.Transitions = append(st.Transitions, tr)
	}
	if len(st.Transitions) == 0 {
		st.Transitions = nil
	}
}

// TransitionSpan is the moment the last transition in a style reaches its
// target: the longest duration-plus-delay. A style whose transitions all run
// instantly reports zero, which a caller reads as "snap, do not animate".
func TransitionSpan(st Style) float64 {
	span := 0.0
	for _, tx := range st.Transitions {
		if end := tx.Delay.MS() + tx.Duration.MS(); end > span {
			span = end
		}
	}
	return span
}

// parsePropertyList reads a transition-property value: a comma-separated
// series of canonical names, or "none" meaning nothing transitions. The
// boolean reports the none flavour; a real list reports false.
func parsePropertyList(raw string) ([]string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if raw == "none" {
		return nil, true
	}
	var out []string
	for _, p := range splitFields(raw, ',') {
		p = strings.ToLower(strings.TrimSpace(p))
		if p == "" {
			return nil, false
		}
		out = append(out, p)
	}
	return out, false
}

// parseTimeList reads a comma-separated list of CSS times.
func parseTimeList(raw string) ([]Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []Time
	for _, p := range splitFields(raw, ',') {
		t, ok := ParseTime(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, t)
	}
	return out, true
}

// parseTimingList reads a comma-separated list of easings.
func parseTimingList(raw string) ([]Timing, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []Timing
	for _, p := range splitFields(raw, ',') {
		t, ok := parseTiming(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, t)
	}
	return out, true
}
