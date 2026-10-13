package css

// parseClear reads the sides a box asks to be pushed past.
func parseClear(raw string) (uint8, bool) {
	switch keyword(raw) {
	case "none":
		return ClearNone, true
	case "left":
		return ClearLeft, true
	case "right":
		return ClearRight, true
	case "both":
		return ClearBoth, true
	}
	return 0, false
}
