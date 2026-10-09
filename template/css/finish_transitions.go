package css

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
