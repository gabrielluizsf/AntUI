package css

func applyTransitionDecl(e *declEnv) []string {
	st, raw, note := e.st, e.raw, e.note
	switch e.prop {
	case "transition":
		if list, ok := parseTransitions(raw); ok {
			st.TransitionProps, st.TransitionDurs, st.TransitionTims, st.TransitionDels = transitionParts(list)
			st.TransitionNone = len(list) == 0
			note()
		}
	case "transition-property":
		if list, none := parsePropertyList(raw); list != nil || none {
			st.TransitionProps = list
			st.TransitionNone = none
			note()
		}
	case "transition-duration":
		if list, ok := parseTimeList(raw); ok {
			st.TransitionDurs = list
			note()
		}
	case "transition-timing-function":
		if list, ok := parseTimingList(raw); ok {
			st.TransitionTims = list
			note()
		}
	case "transition-delay":
		if list, ok := parseTimeList(raw); ok {
			st.TransitionDels = list
			note()
		}
	}
	return nil
}
