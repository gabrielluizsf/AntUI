package css

// parseJustifyContent turns the six justify-content keywords into a
// Justify* constant. The bare start/end keywords behave like their
// flex- counterparts.
func parseJustifyContent(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch":
		return JustifyFlexStart, true
	case "flex-start", "start", "left":
		return JustifyFlexStart, true
	case "flex-end", "end", "right":
		return JustifyFlexEnd, true
	case "center":
		return JustifyCenter, true
	case "space-between":
		return JustifySpaceBetween, true
	case "space-around":
		return JustifySpaceAround, true
	case "space-evenly":
		return JustifySpaceEvenly, true
	}
	return 0, false
}
