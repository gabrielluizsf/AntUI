package css

func applyAnimationDecl(e *declEnv) []string {
	st, raw, note := e.st, e.raw, e.note
	switch e.prop {
	case "animation":
		if list, ok := parseAnimations(raw); ok {
			st.AnimationNames, st.AnimationDurs, st.AnimationTims, st.AnimationDels, st.AnimationIters, st.AnimationDirs, st.AnimationFills = animationParts(list)
			note()
		}
	case "animation-name":
		if list, ok := parseNameList(raw); ok {
			st.AnimationNames = list
			note()
		}
	case "animation-duration":
		if list, ok := parseTimeList(raw); ok {
			st.AnimationDurs = list
			note()
		}
	case "animation-timing-function":
		if list, ok := parseTimingList(raw); ok {
			st.AnimationTims = list
			note()
		}
	case "animation-delay":
		if list, ok := parseTimeList(raw); ok {
			st.AnimationDels = list
			note()
		}
	case "animation-iteration-count":
		if list, ok := parseIterationList(raw); ok {
			st.AnimationIters = list
			note()
		}
	case "animation-direction":
		if list, ok := parseDirectionList(raw); ok {
			st.AnimationDirs = list
			note()
		}
	case "animation-fill-mode":
		if list, ok := parseFillList(raw); ok {
			st.AnimationFills = list
			note()
		}
	}
	return nil
}
