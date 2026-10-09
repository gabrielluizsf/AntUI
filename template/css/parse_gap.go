package css

// parseGap reads a gap shorthand: one length for both axes, or a row gap
// followed by a column gap. A negative length or an unreadable token drops
// the whole declaration.
func parseGap(raw string, ctx Units) (row, col Length, ok bool) {
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 2 {
		return row, col, false
	}
	vals := make([]Length, len(parts))
	for i, p := range parts {
		if p == "normal" {
			vals[i] = Zero()
			continue
		}
		l, err := parseLengthAt(p, ctx)
		if err != nil || l.Auto() || l.None() || l.value < 0 {
			return row, col, false
		}
		vals[i] = l
	}
	if len(vals) == 1 {
		return vals[0], vals[0], true
	}
	return vals[0], vals[1], true
}
