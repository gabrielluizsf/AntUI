package css

// parseOverflowWrap reads the overflow-wrap (or word-wrap) keyword.
func parseOverflowWrap(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return OverflowWrapNormal, true
	case "break-word":
		return OverflowWrapBreakWord, true
	case "anywhere":
		return OverflowWrapAnywhere, true
	}
	return 0, false
}
