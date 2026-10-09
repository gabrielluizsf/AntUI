package css

// parseAlignItems turns the align-items keywords into an Align* constant;
// normal is stretch, the initial value.
func parseAlignItems(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch":
		return AlignStretch, true
	case "flex-start", "start", "self-start":
		return AlignFlexStart, true
	case "flex-end", "end", "self-end":
		return AlignFlexEnd, true
	case "center":
		return AlignCenter, true
	case "baseline":
		return AlignBaseline, true
	}
	return 0, false
}
