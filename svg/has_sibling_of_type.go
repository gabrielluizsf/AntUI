package svg

// hasSiblingOfType says whether the element at n has a sibling of its own name
// written after it (next) or before it. An `-of-type` counts the siblings that
// share the element's name and passes over the rest.
func hasSiblingOfType(n *selNode, next bool) bool {
	for s := siblingStep(n, next); s != nil; s = siblingStep(s, next) {
		if s.el.Name == n.el.Name {
			return true
		}
	}
	return false
}

// siblingStep is the next sibling towards the end (next) or the beginning.
func siblingStep(n *selNode, next bool) *selNode {
	if next {
		return n.next
	}
	return n.prev
}
