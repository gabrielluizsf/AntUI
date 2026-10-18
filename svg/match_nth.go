package svg

// matchNth says whether the element stands where its `:nth-…` names it: its
// position among the siblings the variant counts — the same-name ones for
// `-of-type`, those an `of S` list names where one is written — read from the
// first or the last according to the kind. An element the `of S` list does not
// name is not among the siblings counted at all, so it never matches.
func matchNth(s simpleSelector, n *selNode) bool {
	if len(s.nth.of) > 0 && !matchSub(s.nth.of, n) {
		return false
	}
	ofType := s.kind == kindNthOfType || s.kind == kindNthLastOfType
	last := s.kind == kindNthLastChild || s.kind == kindNthLastOfType
	return matchesNth(s.nth.a, s.nth.b, nthPosition(s, n, ofType, last))
}

// nthPosition counts over the siblings a kind counts, from the element towards
// the first (or towards the last), and answers its 1-based position among them.
func nthPosition(s simpleSelector, n *selNode, ofType, last bool) int {
	pos := 0
	for x := n; x != nil; x = siblingStep(x, last) {
		if ofType && x.el.Name != n.el.Name {
			continue
		}
		if len(s.nth.of) > 0 && !matchSub(s.nth.of, x) {
			continue
		}
		pos++
	}
	return pos
}
