package css

func applyBorderDecl(e *declEnv) []string {
	st, set, raw, note := e.st, e.set, e.raw, e.note
	switch e.prop {
	case "border":
		applyBorder(st, set, allSides, raw, note)
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
		applyBorder(st, set, oneSide(SideTop), raw, note)
	case "border-right":
		applyBorder(st, set, oneSide(SideRight), raw, note)
	case "border-bottom":
		applyBorder(st, set, oneSide(SideBottom), raw, note)
	case "border-left":
		applyBorder(st, set, oneSide(SideLeft), raw, note)
	case "border-radius":
		if rx, ry, ok := parseRadiiXY(raw); ok {
			st.Radius = rx
			st.RadiusY = ry
			note()
		}
	}
	return nil
}
