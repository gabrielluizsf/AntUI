package css

func applyFlowDecl(e *declEnv) []string {
	st, set, raw, prop, note := e.st, e.set, e.raw, e.prop, e.note
	switch prop {
	case "break-inside", "page-break-inside":
		if v, ok := parseBreakInside(raw); ok {
			st.BreakInside = v
			note()
			set["break-inside"] = true
		}
	case "float":
		if v, ok := parseFloat(raw); ok {
			st.Float = v
			note()
			if v != FloatNone {
				return []string{fmtErrf("%s is read but not applied: this template's flow has no line box to float inside", prop).Error()}
			}
		}
	case "clear":
		if v, ok := parseClear(raw); ok {
			st.Clear = v
			note()
			if v != ClearNone {
				return []string{fmtErrf("%s is read but not applied: nothing here floats, so there is nothing to be pushed past", prop).Error()}
			}
		}
	}
	return nil
}
