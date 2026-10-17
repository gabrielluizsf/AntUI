package svg

// matchSimple says whether one condition holds for the element at node: its
// type, a class it carries, an id it wrote, the universal selector (which
// always does), an attribute condition, or a pseudo-class about where it
// stands among its siblings.
func matchSimple(s simpleSelector, n *selNode) bool {
	switch s.kind {
	case kindUniversal:
		return true
	case kindType:
		return n.el.Name == s.name
	case kindClass:
		return hasClass(n.el.attr("class"), s.name)
	case kindID:
		return n.el.attr("id") == s.name
	case kindAttribute:
		return matchAttr(n.el, s)
	case kindRoot:
		return n.parent == nil
	case kindFirstChild:
		return n.prev == nil
	case kindLastChild:
		return n.next == nil
	case kindOnlyChild:
		return n.prev == nil && n.next == nil
	}
	return false
}
