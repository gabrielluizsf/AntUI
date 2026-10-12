package svg

// readOffset reads a `feOffset`: the dx and dy it moves the picture by.
func (p *filterPrim) readOffset(e *element, warn func(string, ...any)) bool {
	for _, at := range []struct {
		attr string
		slot *float64
	}{{"dx", &p.dx}, {"dy", &p.dy}} {
		if raw := e.attr(at.attr); raw != "" {
			v, ok := parseLength(raw)
			if !ok {
				warn("the feOffset %s %q is not a length this package can read, so the primitive is left out", at.attr, raw)
				return false
			}
			*at.slot = v
		}
	}
	return true
}
