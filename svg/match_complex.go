package svg

// matchComplex says whether a selector stands on the element at node: the
// rightmost compound matches it, and every compound to its left is reachable
// through the combinators between them.
func matchComplex(cs complexSelector, node *selNode) bool {
	last := len(cs.compounds) - 1
	if !matchCompound(cs.compounds[last], node) {
		return false
	}
	return matchLeft(cs, last, node)
}

// matchLeft walks the compounds to the left of the one already matched,
// following each combinator: a descendant may be any element above, a child
// only the one directly above, and the two sibling combinators the elements
// written before at the same level.
func matchLeft(cs complexSelector, i int, node *selNode) bool {
	if i == 0 {
		return true
	}
	comb, comp := cs.combs[i-1], cs.compounds[i-1]
	return stepBack(comb, node, func(x *selNode) bool {
		return matchCompound(comp, x) && matchLeft(cs, i-1, x)
	})
}
