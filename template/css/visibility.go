package css

// Visibility controls whether a box is painted; a hidden box keeps its slot.
const (
	VisibilityVisible uint8 = iota
	VisibilityHidden
	VisibilityCollapse
)

func parseVisibility(raw string) (uint8, bool) {
	switch raw {
	case "visible":
		return VisibilityVisible, true
	case "hidden":
		return VisibilityHidden, true
	case "collapse":
		return VisibilityCollapse, true
	}
	return 0, false
}
