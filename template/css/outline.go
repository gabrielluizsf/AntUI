package css

// outlineStyle reads one outline-style keyword into the Border* constants an
// outline shares with a border. "none" is the initial value.
func outlineStyle(raw string) (uint8, bool) {
	switch raw {
	case "none":
		return BorderNone, true
	case "solid":
		return BorderSolid, true
	case "dashed":
		return BorderDashed, true
	case "dotted":
		return BorderDotted, true
	case "double":
		return BorderDouble, true
	}
	return BorderNone, false
}

// applyOutline folds the outline shorthand into the style: a width, a style
// and a colour in any order, or the keyword none to clear it. It reports
// whether anything was recognised.
func applyOutline(st *Style, raw string, ctx Units) bool {
	if raw == "none" {
		st.OutlineWidth = 0
		st.OutlineStyle = BorderNone
		return true
	}
	matched := false
	for _, tk := range splitTokens(raw) {
		if s, ok := outlineStyle(tk); ok {
			st.OutlineStyle = s
			matched = true
			continue
		}
		if c, err := ParseColor(tk); err == nil {
			st.OutlineColor = c
			matched = true
			continue
		}
		l, err := parseLengthAt(tk, ctx)
		if err != nil || l.IsPct() || l.Auto() || l.None() {
			return false
		}
		st.OutlineWidth = l.Resolve(ctx)
		matched = true
	}
	return matched
}
