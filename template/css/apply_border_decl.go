package css

func applyBorderDecl(e *declEnv) []string {
	st, set, raw, note := e.st, e.set, e.raw, e.note
	switch e.prop {
	case "border":
		applyBorder(st, set, [4]int{0, 1, 2, 3}, raw, note)
	case "border-width":
		if v, ok := parseBorderWidths(raw); ok {
			st.BorderWidth = v
			note()
		}
	case "border-style":
		if v, ok := parseBorderStyles(raw); ok {
			st.BoxStyle = v
			note()
		}
	case "border-color":
		if v, ok := parseBorderColors(raw); ok {
			st.BoxColor = v
			note()
		}
	case "border-top":
		applyBorder(st, set, [4]int{0, -1, -1, -1}, raw, note)
	case "border-right":
		applyBorder(st, set, [4]int{-1, 1, -1, -1}, raw, note)
	case "border-bottom":
		applyBorder(st, set, [4]int{-1, -1, 2, -1}, raw, note)
	case "border-left":
		applyBorder(st, set, [4]int{-1, -1, -1, 3}, raw, note)
	case "border-radius":
		if rx, ry, ok := parseRadiiXY(raw); ok {
			st.Radius = rx
			st.RadiusY = ry
			note()
		}
	}
	return nil
}
