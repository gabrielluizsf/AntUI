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
	switch comb {
	case child:
		return node.parent != nil && matchCompound(comp, node.parent) && matchLeft(cs, i-1, node.parent)
	case adjacent:
		return node.prev != nil && matchCompound(comp, node.prev) && matchLeft(cs, i-1, node.prev)
	case sibling:
		for s := node.prev; s != nil; s = s.prev {
			if matchCompound(comp, s) && matchLeft(cs, i-1, s) {
				return true
			}
		}
	case descendant:
		for a := node.parent; a != nil; a = a.parent {
			if matchCompound(comp, a) && matchLeft(cs, i-1, a) {
				return true
			}
		}
	}
	return false
}
