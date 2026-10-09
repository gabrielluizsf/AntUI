package css

func applyTextDecl(e *declEnv) []string {
	st, set, raw, ctx, note := e.st, e.set, e.raw, e.ctx, e.note
	switch e.prop {
	case "text-align":
		switch raw {
		case "left":
			st.TextAlign = 0
			note()
		case "center":
			st.TextAlign = 1
			note()
		case "right":
			st.TextAlign = 2
			note()
		case "justify":
			st.TextAlign = 3
			note()
		}
	case "text-transform":
		if v, ok := parseTextTransform(raw); ok {
			st.TextTransform = v
			note()
		}
	case "text-decoration", "text-decoration-line":
		if bits, ok := parseTextDecoration(raw); ok {
			st.TextDecoration = bits
			note()
		}
	case "white-space":
		if v, ok := parseWhiteSpace(raw); ok {
			st.WhiteSpace = v
			note()
		}
	case "overflow-wrap", "word-wrap":
		if v, ok := parseOverflowWrap(raw); ok {
			st.OverflowWrap = v
			note()
			set["overflow-wrap"] = true
		}
	case "text-overflow":
		if v, ok := parseTextOverflow(raw); ok {
			st.TextOverflow = v
			note()
		}
	case "vertical-align":
		if v, l, ok := parseVerticalAlign(raw, *ctx); ok {
			st.VerticalAlign = v
			st.BaselineShift = l
			note()
		}
	}
	return nil
}
