package svg

import (
	"fmt"
	"strings"
)

// nthKind is the kind each `:nth-…` parses into.
var nthKind = map[string]simpleKind{
	"nth-child":        kindNthChild,
	"nth-last-child":   kindNthLastChild,
	"nth-of-type":      kindNthOfType,
	"nth-last-of-type": kindNthLastOfType,
}

// appendNth reads an `:nth-…(…)`: an `An+B` expression and, for the two
// positional ones, the `of S` list that narrows the siblings counted.
func appendNth(c *compoundSelector, s, name string, next, at int, ns namespace) (int, error) {
	if next >= len(s) || s[next] != '(' {
		return at, fmt.Errorf(":%s needs an argument in parentheses", name)
	}
	inner, after, err := readParen(s, next)
	if err != nil {
		return at, err
	}
	form, of := splitNthOf(inner)
	nth, ok := parseNth(form)
	if !ok {
		return at, fmt.Errorf("%q is not an An+B expression", strings.TrimSpace(form))
	}
	kind := nthKind[name]
	if of != "" {
		if kind == kindNthOfType || kind == kindNthLastOfType {
			return at, fmt.Errorf(":%s takes no `of` list", name)
		}
		sub, err := parseSelectorList(of, ns)
		if err != nil {
			return at, fmt.Errorf(":%s: %w", name, err)
		}
		if len(sub) == 0 {
			return at, fmt.Errorf(":%s has no selector after `of`", name)
		}
		nth.of = sub
	}
	c.simples = append(c.simples, simpleSelector{kind: kind, nth: nth})
	return after, nil
}
