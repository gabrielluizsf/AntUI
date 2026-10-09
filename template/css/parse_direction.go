package css

// parseDirection reads one direction keyword.
func parseDirection(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return AnimNormal, true
	case "reverse":
		return AnimReverse, true
	case "alternate":
		return AnimAlternate, true
	case "alternate-reverse":
		return AnimAlternateReverse, true
	}
	return 0, false
}
