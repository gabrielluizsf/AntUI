package svg

// stepBack walks one combinator back from a node that stands on its right and
// tests what is on the left: the parent for a child combinator, the previous
// sibling for a `+`, any earlier sibling for a `~`, any ancestor for a
// descendant.
func stepBack(comb combinator, from *selNode, test func(*selNode) bool) bool {
	switch comb {
	case child:
		return from.parent != nil && test(from.parent)
	case adjacent:
		return from.prev != nil && test(from.prev)
	case sibling:
		for s := from.prev; s != nil; s = s.prev {
			if test(s) {
				return true
			}
		}
	case descendant:
		for a := from.parent; a != nil; a = a.parent {
			if test(a) {
				return true
			}
		}
	}
	return false
}
