package svg

import "fmt"

// appendUniversal reads the universal selector, `*` or the `*|` that names any
// namespace.
func appendUniversal(c *compoundSelector, s string, i int, ns namespace) (int, error) {
	if i+1 < len(s) && s[i+1] == '|' && (i+2 >= len(s) || s[i+2] != '|') {
		return appendLocalType(c, s, i+2, "*", ns)
	}
	c.simples = append(c.simples, simpleSelector{kind: kindUniversal})
	return i + 1, nil
}

// appendNoNamespaceType reads a type selector with an empty namespace prefix,
// `|rect`. No element a drawing paints is in no namespace — an element with
// none is left out of the tree — so it never matches, and is kept as the valid
// selector that it is.
func appendNoNamespaceType(c *compoundSelector, s string, i int) (int, error) {
	if i < len(s) && s[i] == '*' {
		c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable})
		return i + 1, nil
	}
	name, next := readSelIdent(s, i)
	if name == "" {
		return i, fmt.Errorf("a namespace prefix names no type")
	}
	c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable})
	return next, nil
}
