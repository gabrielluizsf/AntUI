package svg

import "fmt"

// appendPseudo reads one pseudo-class or pseudo-element, `:name` or `::name`,
// with the parenthesised argument a functional one is given. i is the index of
// the first `:`, and ns is the namespaces in scope where the sheet stands,
// which an argument carrying a selector reads through.
//
// A pseudo-element never matches, because a still drawing has no generated box
// for one; a pseudo-class that asks about a state a drawing cannot be in —
// `:hover`, `:focus` — never matches either. Both are read and kept in the
// selector rather than refused, so a rule written with them stands aside and
// the drawing is painted as if the rule had not named it. A name that is not a
// CSS pseudo at all is a different thing: the selector is not one, and the
// rule is left out.
func appendPseudo(c *compoundSelector, s string, i int, ns namespace) (int, error) {
	start := i + 1
	pe := start < len(s) && s[start] == ':'
	if pe {
		start++
	}
	name, next := readSelIdent(s, start)
	if name == "" {
		return i, fmt.Errorf("a pseudo-selector is missing its name")
	}
	if pe {
		if next < len(s) && s[next] == '(' {
			var err error
			_, next, err = readParen(s, next)
			if err != nil {
				return i, err
			}
		}
		c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable, pe: true})
		return next, nil
	}
	return appendPseudoClass(c, s, name, next, i, ns)
}
