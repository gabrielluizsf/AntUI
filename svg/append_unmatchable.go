package svg

// appendUnmatchable reads the argument a valid but unanswerable functional
// pseudo-class may carry — `:lang(en)`, `:host(.x)` — and adds the condition
// that never holds.
func appendUnmatchable(c *compoundSelector, s string, next, at int) (int, error) {
	next, err := skipParenArg(s, next, at)
	if err != nil {
		return at, err
	}
	c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable})
	return next, nil
}

// appendUnmatchableElement is [appendUnmatchable] for a pseudo-element, whose
// argument a functional one may carry, `::part(x)`.
func appendUnmatchableElement(c *compoundSelector, s string, next, at int) (int, error) {
	next, err := skipParenArg(s, next, at)
	if err != nil {
		return at, err
	}
	c.simples = append(c.simples, simpleSelector{kind: kindUnmatchable, pe: true})
	return next, nil
}

// skipParenArg passes the parenthesised argument at next, where there is one,
// and answers where it ends.
func skipParenArg(s string, next, at int) (int, error) {
	if next >= len(s) || s[next] != '(' {
		return next, nil
	}
	_, after, err := readParen(s, next)
	if err != nil {
		return at, err
	}
	return after, nil
}
