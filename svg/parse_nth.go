package svg

import (
	"strconv"
	"strings"
)

// parseNth reads an `An+B` expression: `odd`, `even`, a plain integer, or a
// form naming `n`, with optional signs and the whitespace around them allowed.
// It answers the two coefficients and whether the form was one at all.
func parseNth(s string) (nthExpr, bool) {
	s = strings.ToLower(stripNthSpace(s))
	switch s {
	case "odd":
		return nthExpr{a: 2, b: 1}, true
	case "even":
		return nthExpr{a: 2, b: 0}, true
	}
	i := strings.IndexByte(s, 'n')
	if i < 0 {
		b, err := strconv.Atoi(s)
		if err != nil {
			return nthExpr{}, false
		}
		return nthExpr{b: b}, true
	}
	a, ok := nthCoeff(s[:i])
	if !ok {
		return nthExpr{}, false
	}
	b := 0
	if rest := s[i+1:]; rest != "" {
		if rest[0] != '+' && rest[0] != '-' {
			return nthExpr{}, false
		}
		v, err := strconv.Atoi(rest)
		if err != nil {
			return nthExpr{}, false
		}
		b = v
	}
	return nthExpr{a: a, b: b}, true
}
