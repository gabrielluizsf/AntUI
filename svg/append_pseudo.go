package svg

import "fmt"

// appendPseudo reads a `:name` pseudo-class. Only the ones a drawing read once
// can answer are read — where an element stands among its siblings, and the
// root — and anything else is an error, so a rule that asks for a state a
// still drawing has none of stands aside rather than matching what it did not
// mean to.
func appendPseudo(c *compoundSelector, s string, i int) (int, error) {
	name, next := readSelIdent(s, i)
	switch name {
	case "root":
		c.simples = append(c.simples, simpleSelector{kind: kindRoot})
	case "first-child":
		c.simples = append(c.simples, simpleSelector{kind: kindFirstChild})
	case "last-child":
		c.simples = append(c.simples, simpleSelector{kind: kindLastChild})
	case "only-child":
		c.simples = append(c.simples, simpleSelector{kind: kindOnlyChild})
	default:
		return i, fmt.Errorf("the pseudo-class :%s is not one a drawing can answer", name)
	}
	return next, nil
}
