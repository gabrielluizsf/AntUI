package css

func applyFlexDecl(e *declEnv) []string {
	st, raw, note := e.st, e.raw, e.note
	switch e.prop {
	case "flex-direction":
		if v, ok := parseFlexDirection(raw); ok {
			st.FlexDirection = v
			note()
		}
	case "flex-wrap":
		if v, ok := parseFlexWrap(raw); ok {
			st.FlexWrap = v
			note()
		}
	case "flex-flow":
		words := splitWords(raw)
		if len(words) == 0 || len(words) > 2 {
			break
		}
		ok := true
		for _, w := range words {
			if v, good := parseFlexDirection(w); good {
				st.FlexDirection = v
				continue
			}
			if v, good := parseFlexWrap(w); good {
				st.FlexWrap = v
				continue
			}
			ok = false
		}
		if ok {
			note()
		}
	}
	return nil
}
