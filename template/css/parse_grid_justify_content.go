package css

// parseGridJustifyContent reads the grid alignment keywords for the rail content.
func parseGridJustifyContent(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch", "left", "start", "flex-start":
		if raw == "normal" || raw == "stretch" {
			return ContentStretch, true
		}
		return ContentFlexStart, true
	case "right", "end", "flex-end":
		return ContentFlexEnd, true
	case "center":
		return ContentCenter, true
	case "space-between":
		return ContentSpaceBetween, true
	case "space-around":
		return ContentSpaceAround, true
	case "space-evenly":
		return ContentSpaceEvenly, true
	}
	return 0, false
}
