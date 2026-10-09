package css

func applyFlexItemDecl(e *declEnv) []string {
	st, set, raw, ctx, note := e.st, e.set, e.raw, e.ctx, e.note
	switch e.prop {
	case "order":
		if v, ok := parseOrder(raw); ok {
			st.Order = v
			note()
		}
	case "flex":
		if g, s, b, ok := parseFlex(raw, *ctx); ok {
			st.FlexGrow, st.FlexShrink, st.FlexBasis = g, s, b
			note()
			set["flex-grow"] = true
			set["flex-shrink"] = true
			set["flex-basis"] = true
		}
	case "flex-grow":
		if v, ok := parseFlexNumber(raw); ok {
			st.FlexGrow = v
			note()
		}
	case "flex-shrink":
		if v, ok := parseFlexNumber(raw); ok {
			st.FlexShrink = v
			note()
		}
	case "flex-basis":
		if v, ok := parseFlexBasis(raw, *ctx); ok {
			st.FlexBasis = v
			note()
		}
	}
	return nil
}
