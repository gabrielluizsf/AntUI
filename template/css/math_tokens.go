package css

// mathTokens splits a formula into operators, parentheses, commas and
// numbers-with-units. "calc(" and "min(" stay whole so the parser matches the
// function name before its "(".
type mathTokens struct {
	s    string
	toks []string
	i    int
}
