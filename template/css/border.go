package css

// Border styles, matching the values stored in Style.BoxStyle.
const (
	BorderNone uint8 = iota
	BorderSolid
	BorderDashed
	BorderDotted
	BorderDouble
)

func parseBorderStyle(raw string) (uint8, bool) {
	switch raw {
	case "solid":
		return BorderSolid, true
	case "dashed":
		return BorderDashed, true
	case "dotted":
		return BorderDotted, true
	case "double":
		return BorderDouble, true
	case "groove", "ridge", "inset", "outset":
		// Drawn as a solid ring: a flat canvas has no bevel to shade.
		return BorderSolid, true
	case "none", "hidden":
		return BorderNone, true
	}
	return 0, false
}
