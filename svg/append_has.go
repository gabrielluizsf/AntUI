package svg

import "fmt"

// appendHas reads `:has(…)`: a list of relative selectors, each read from the
// element outward — `:has(> rect)`, `:has(+ p .x)` — through the namespaces in
// scope where the sheet stands.
func appendHas(c *compoundSelector, s string, next, at int, ns namespace) (int, error) {
	if next >= len(s) || s[next] != '(' {
		return at, fmt.Errorf(":has needs a selector list in parentheses")
	}
	inner, after, err := readParen(s, next)
	if err != nil {
		return at, err
	}
	rel, err := parseRelativeList(inner, ns)
	if err != nil {
		return at, fmt.Errorf(":has: %w", err)
	}
	if len(rel) == 0 {
		return at, fmt.Errorf(":has names no selector")
	}
	c.simples = append(c.simples, simpleSelector{kind: kindHas, has: rel})
	return after, nil
}
