package css

// transformLength is a length a transform accepts: a calc(), a percentage or
// a fixed unit surface, but never the auto/none keywords.
func transformLength(s string, ctx Units) (Length, bool) {
	l, err := parseLengthAt(s, ctx)
	if err != nil || l.Auto() || l.None() {
		return Length{}, false
	}
	return l, true
}
