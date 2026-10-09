package css

import "strings"

// parseMediaQuery parses one comma-separated media query into its alternative
// width windows, or-separated groups each becoming a query of their own.
func parseMediaQuery(branch string) ([]MediaQuery, []string) {
	var out []MediaQuery
	var warns []string
	negated := false
	cur := MediaQuery{}
	flush := func() {
		if !cur.constrained() {
			// A negation of nothing we can evaluate is left unconstrained so
			// the rule still applies.
			cur.Negated = false
		}
		out = append(out, cur)
	}
	for _, tok := range splitQueryTokens(branch) {
		t := strings.ToLower(strings.TrimSpace(tok))
		if t == "" {
			continue
		}
		switch t {
		case "not":
			negated = true
			cur.Negated = true
		case "only":
			// A media-type qualifier; the type itself is beyond the canvas.
		case "and":
			// Conjunction: keep filling the current query.
		case "or":
			flush()
			cur = MediaQuery{Negated: negated}
		default:
			warns = append(warns, readMediaToken(&cur, t)...)
		}
	}
	flush()
	return out, warns
}
