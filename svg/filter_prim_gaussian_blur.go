package svg

// readGaussianBlur reads a `feGaussianBlur`: the standard deviation, which may
// be one number for both axes or one for each, and must not be negative.
func (p *filterPrim) readGaussianBlur(e *element, warn func(string, ...any)) bool {
	raw := e.attr("stdDeviation")
	if raw == "" {
		return true
	}
	ns, ok := numbers(raw)
	if !ok || len(ns) == 0 {
		warn("the feGaussianBlur stdDeviation %q is not a number this package can read, so the primitive is left out", raw)
		return false
	}
	p.stdDev, p.stdDevY = ns[0], ns[0]
	if p.stdDev < 0 {
		warn("the feGaussianBlur stdDeviation %q is negative, so the primitive is left out", raw)
		return false
	}
	if len(ns) > 1 {
		if ns[1] < 0 {
			warn("the feGaussianBlur stdDeviation %q is negative, so the primitive is left out", raw)
			return false
		}
		p.stdDevY = ns[1]
	}
	return true
}
