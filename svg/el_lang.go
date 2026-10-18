package svg

// elLang answers the language of the element at node: the nearest `xml:lang`
// or `lang` above it that gives one, since language travels down from whoever
// declared it. An empty value says no language at all, which stops the search
// rather than falling through to an ancestor.
func elLang(n *selNode) string {
	for m := n; m != nil; m = m.parent {
		if v, ok := m.el.Attr["xml:lang"]; ok {
			return v
		}
		if v, ok := m.el.Attr["lang"]; ok {
			return v
		}
	}
	return ""
}
