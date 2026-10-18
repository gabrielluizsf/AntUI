package svg

import (
	"fmt"
	"strings"
)

// appendLang reads `:lang(…)` — one or more language ranges an element's
// writing may stand in, separated by commas and each an identifier or a
// quoted string, as CSS writes them.
func appendLang(c *compoundSelector, s string, next, at int) (int, error) {
	if next >= len(s) || s[next] != '(' {
		return at, fmt.Errorf(":lang needs an argument in parentheses")
	}
	inner, after, err := readParen(s, next)
	if err != nil {
		return at, err
	}
	var langs []string
	for _, part := range splitTopComma(inner) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if n := len(part); n >= 2 && (part[0] == '"' || part[0] == '\'') {
			part = part[1 : n-1]
		}
		langs = append(langs, part)
	}
	if len(langs) == 0 {
		return at, fmt.Errorf(":lang names no language")
	}
	c.simples = append(c.simples, simpleSelector{kind: kindLang, lang: langs})
	return after, nil
}
