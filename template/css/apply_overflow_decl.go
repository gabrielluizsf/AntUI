package css

func applyOverflowDecl(e *declEnv) []string {
	st, set, raw, note := e.st, e.set, e.raw, e.note
	switch e.prop {
	case "overflow":
		if v, ok := parseOverflow(raw); ok {
			st.Overflow = [2]uint8{v, v}
			note()
			set["overflow-x"] = true
			set["overflow-y"] = true
		}
	case "overflow-x":
		if v, ok := parseOverflow(raw); ok {
			st.Overflow[0] = v
			note()
			set["overflow"] = true
		}
	case "overflow-y":
		if v, ok := parseOverflow(raw); ok {
			st.Overflow[1] = v
			note()
			set["overflow"] = true
		}
	}
	return nil
}
