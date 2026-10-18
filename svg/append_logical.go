package svg

import "fmt"

// appendLogical reads `:not(…)`, `:is(…)` or `:where(…)`: a selector list
// whose members stand on the element the compound is matched against, read
// through the namespaces in scope where the sheet stands.
func appendLogical(c *compoundSelector, s, name string, next, at int, ns namespace) (int, error) {
	if next >= len(s) || s[next] != '(' {
		return at, fmt.Errorf(":%s needs a selector list in parentheses", name)
	}
	inner, after, err := readParen(s, next)
	if err != nil {
		return at, err
	}
	sub, err := parseSelectorList(inner, ns)
	if err != nil {
		return at, fmt.Errorf(":%s: %w", name, err)
	}
	if len(sub) == 0 {
		return at, fmt.Errorf(":%s names no selector", name)
	}
	kind := kindIs
	switch name {
	case "not":
		kind = kindNot
	case "where":
		kind = kindWhere
	}
	c.simples = append(c.simples, simpleSelector{kind: kind, sub: sub})
	return after, nil
}
