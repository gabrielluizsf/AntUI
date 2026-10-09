package css

// isStopHint reports whether a stop is a bare offset with no colour, the
// colour-hint form the engine skips.
func isStopHint(p string, ctx Units) bool {
	if _, err := parseStopOffset(p, ctx); err != nil {
		return false
	}
	if _, err := ParseColor(p); err == nil {
		return false
	}
	return true
}
