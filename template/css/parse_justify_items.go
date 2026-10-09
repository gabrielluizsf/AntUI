package css

// parseJustifyItems reads the justify-items keywords for the cell contents.
func parseJustifyItems(raw string) (uint8, bool) {
	switch raw {
	case "normal", "stretch", "left":
		if raw == "normal" || raw == "stretch" {
			return AlignStretch, true
		}
		return AlignFlexStart, true
	case "start", "self-start":
		return AlignFlexStart, true
	case "right", "end", "self-end":
		return AlignFlexEnd, true
	case "center":
		return AlignCenter, true
	}
	return 0, false
}
