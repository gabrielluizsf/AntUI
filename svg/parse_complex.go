package svg

import "fmt"

// parseComplex reads one selector — compound selectors joined by combinators —
// into the chain it names, against the namespaces in scope where the sheet
// stands.
func parseComplex(s string, ns namespace) (complexSelector, error) {
	var cs complexSelector
	i := skipSelSpace(s, 0)
	for {
		comp, next, err := parseCompound(s, i, ns)
		if err != nil {
			return cs, err
		}
		cs.compounds = append(cs.compounds, comp)
		j := skipSelSpace(s, next)
		if j >= len(s) {
			return cs, nil
		}
		comb := descendant
		if c := s[j]; c == '>' || c == '+' || c == '~' {
			comb = combinatorOf(c)
			j = skipSelSpace(s, j+1)
		} else if j == next {
			return cs, fmt.Errorf("two conditions stand together with nothing between them at %q", s[j:])
		}
		cs.combs = append(cs.combs, comb)
		i = j
	}
}

// combinatorOf is the combinator a selector wrote between two compounds.
func combinatorOf(c byte) combinator {
	switch c {
	case '>':
		return child
	case '+':
		return adjacent
	}
	return sibling
}

// skipSelSpace passes the whitespace between two pieces of a selector.
func skipSelSpace(s string, i int) int {
	for i < len(s) && isSpace(s[i]) {
		i++
	}
	return i
}
