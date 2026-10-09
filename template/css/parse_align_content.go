package css

// parseAlignContent turns the align-content keywords into a Content*
// constant; normal behaves like stretch.
func parseAlignContent(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch":
		return ContentStretch, true
	case "flex-start", "start":
		return ContentFlexStart, true
	case "flex-end", "end":
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
