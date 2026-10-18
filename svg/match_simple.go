package svg

// matchSimple says whether one condition holds for the element at node: its
// type, a class it carries, an id it wrote, the universal selector (which
// always does), an attribute condition, the language or direction it is
// written in, or a pseudo-class. A pseudo-class a still drawing cannot answer
// never holds, so a rule that names one stands aside instead of matching what
// it did not mean.
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
	case kindRoot, kindScope:
		return n.parent == nil
	case kindEmpty:
		return len(n.el.Kids) == 0 && len(n.el.Parts) == 0
	case kindFirstChild:
		return n.prev == nil
	case kindLastChild:
		return n.next == nil
	case kindOnlyChild:
		return n.prev == nil && n.next == nil
	case kindFirstOfType:
		return !hasSiblingOfType(n, false)
	case kindLastOfType:
		return !hasSiblingOfType(n, true)
	case kindOnlyOfType:
		return !hasSiblingOfType(n, false) && !hasSiblingOfType(n, true)
	case kindNthChild, kindNthLastChild, kindNthOfType, kindNthLastOfType:
		return matchNth(s, n)
	case kindNot:
		return !matchSub(s.sub, n)
	case kindIs, kindWhere:
		return matchSub(s.sub, n)
	case kindHas:
		return matchHasList(s.has, n)
	case kindLang:
		return matchLang(elLang(n), s.lang)
	case kindDir:
		return elDir(n) == s.name
	}
	return false
}
