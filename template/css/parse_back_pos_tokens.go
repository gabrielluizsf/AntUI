package css

// parseBackPosTokens folds a token list into one [BackPos]. A component is a
// keyword with the offset that immediately follows it, or a bare length:
// "right 20px", "top", "center", "10% 20px". Keywords lock their component
// onto their own axis; bare lengths and centre fill the free axes in order,
// and an axis left over after both components settle is centred — so "left"
// means left centre and "20px" means a 20px hand-off on a vertically centred
// picture.
func parseBackPosTokens(tokens []string, ctx Units) (BackPos, bool) {
	comps, ok := readBackComps(tokens, ctx)
	if !ok {
		return BackPos{}, false
	}
	var out BackPos
	placed := [2]bool{}
	assign := func(idx int, p BackPosition) bool {
		if placed[idx] {
			return false
		}
		out[idx] = p
		placed[idx] = true
		return true
	}
	// First pass: components whose keyword names an axis go there. A keyword
	// whose own axis is taken falls back onto the other; a component that
	// finds both taken is malformed.
	for _, c := range comps {
		if c.axis == 2 {
			continue
		}
		if c.axis == 0 {
			if !assign(0, c.pos) && !assign(1, c.pos) {
				return BackPos{}, false
			}
		} else if !assign(1, c.pos) && !assign(0, c.pos) {
			return BackPos{}, false
		}
	}
	// Second pass: bare offsets and centre fill the free axes in order.
	for _, c := range comps {
		if c.axis != 2 {
			continue
		}
		if !assign(0, c.pos) && !assign(1, c.pos) {
			return BackPos{}, false
		}
	}
	// An axis the values never named is centred.
	if !placed[0] {
		out[0] = BackPosition{Edge: BackPosMiddle, Off: Zero()}
	}
	if !placed[1] {
		out[1] = BackPosition{Edge: BackPosMiddle, Off: Zero()}
	}
	return out, true
}
