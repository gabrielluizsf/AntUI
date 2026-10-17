package svg

import "fmt"

// appendNamed reads the name after a `.` or a `#`, or the type name at the
// start of a compound, and adds it.
func appendNamed(c *compoundSelector, kind simpleKind, s string, i int) (int, error) {
	name, next := readSelIdent(s, i)
	if name == "" {
		return i, fmt.Errorf("%q is missing the name a selector needs", s)
	}
	c.simples = append(c.simples, simpleSelector{kind: kind, name: name})
	return next, nil
}
