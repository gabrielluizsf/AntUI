package css

// parseFlex reads the flex shorthand: none, auto, initial, or up to three
// values holding grow, shrink and basis in any order. A unitless number is
// the next-flex factor (grow first, shrink second); a length or a keyword is
// the basis. What the declaration does not say defaults to grow 1, shrink 1,
// basis 0, so "flex: 2" and "flex: 2 1 0%" agree.
func parseFlex(raw string, ctx Units) (grow, shrink float64, basis Length, ok bool) {
	switch raw {
	case "none":
		return 0, 0, Auto(), true
	case "auto":
		return 1, 1, Auto(), true
	case "initial":
		return 0, 1, Auto(), true
	}
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 3 {
		return 0, 0, Length{}, false
	}
	seenGrow, seenShrink := false, false
	hasBasis := false
	for _, p := range parts {
		if v, num := parseFlexNumber(p); num && !seenShrink {
			if !seenGrow {
				grow, seenGrow = v, true
				continue
			}
			shrink, seenShrink = v, true
			continue
		}
		if b, good := parseFlexBasis(p, ctx); good && !hasBasis {
			basis, hasBasis = b, true
			continue
		}
		return 0, 0, Length{}, false
	}
	if !seenGrow {
		grow = 1
	}
	if !seenShrink {
		shrink = 1
	}
	if !hasBasis {
		basis = Pct(0)
	}
	return grow, shrink, basis, true
}
