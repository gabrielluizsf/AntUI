package svg

// matchHasList says whether any relative selector of a `:has` list is reached
// from the element it stands on.
func matchHasList(list []relativeSelector, scope *selNode) bool {
	for _, rel := range list {
		if matchRelative(rel, scope) {
			return true
		}
	}
	return false
}

// matchRelative says whether a relative selector reaches some node from the
// element it stands on: the rightmost compound stands on a node in reach, and
// the compounds to its left are walkable back to the element through the
// combinators, the lead joining the leftmost compound to the element itself.
func matchRelative(rel relativeSelector, scope *selNode) bool {
	last := len(rel.cs.compounds) - 1
	for _, cand := range relCandidates(scope) {
		if matchRelChain(rel, last, cand, scope) {
			return true
		}
	}
	return false
}

// matchRelChain matches the compounds from i down to the left, the leftmost
// then joined to the element by the lead combinator.
func matchRelChain(rel relativeSelector, i int, node, scope *selNode) bool {
	if !matchCompound(rel.cs.compounds[i], node) {
		return false
	}
	if i == 0 {
		return stepBack(rel.lead, node, func(x *selNode) bool { return x == scope })
	}
	return stepBack(rel.cs.combs[i-1], node, func(x *selNode) bool {
		return matchRelChain(rel, i-1, x, scope)
	})
}
