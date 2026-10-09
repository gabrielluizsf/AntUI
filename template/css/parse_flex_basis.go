package css

// parseFlexBasis reads the flex-basis longhand: a length, a percentage,
// auto, or content. content maps to the none marker so the solver reads it
// as "measure the item's own main size", the same way auto is handled.
func parseFlexBasis(raw string, ctx Units) (Length, bool) {
	if raw == "auto" {
		return Auto(), true
	}
	if raw == "content" || raw == "max-content" || raw == "min-content" || raw == "fit-content" {
		return Length{u: unitNone}, true
	}
	l, err := parseLengthAt(raw, ctx)
	if err != nil {
		return l, false
	}
	return l, true
}
