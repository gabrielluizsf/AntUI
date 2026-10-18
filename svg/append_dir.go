package svg

import (
	"fmt"
	"strings"
)

// appendDir reads `:dir(…)`, whose argument is the direction it asks for —
// `ltr` or `rtl` — which the drawing says on the element's own `dir`.
func appendDir(c *compoundSelector, s string, next, at int) (int, error) {
	if next >= len(s) || s[next] != '(' {
		return at, fmt.Errorf(":dir needs an argument in parentheses")
	}
	inner, after, err := readParen(s, next)
	if err != nil {
		return at, err
	}
	arg := strings.ToLower(strings.TrimSpace(inner))
	if arg != "ltr" && arg != "rtl" {
		return at, fmt.Errorf(":dir is given %q, where ltr or rtl was expected", strings.TrimSpace(inner))
	}
	c.simples = append(c.simples, simpleSelector{kind: kindDir, name: arg})
	return after, nil
}
