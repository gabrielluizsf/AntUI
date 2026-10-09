package css

// EvalMath is the length-formula engine behind calc(), min(), max() and
// clamp(): it turns "calc(100% - 20px)", "min(480px, 50vw)" and friends into
// reference pixels under a measurement context. Every term collapses to
// pixels — percentages against the containing width, viewport and font units
// against the relevant measure — and the arithmetic runs in that space, which
// is exactly how a browser expands a length mix. ok is false when the formula
// does not balance or a term does not parse, and the caller warns or drops.
func EvalMath(raw string, ctx Units) (px float64, ok bool) {
	t := &mathTokens{s: stripComments(raw)}
	p := &mathParser{t: t, ctx: ctx}
	v, err := p.expr()
	if err != nil || !t.empty() {
		return 0, false
	}
	return v, true
}
