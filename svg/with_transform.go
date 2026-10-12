package svg

import "strings"

// withTransform reads a `transform`. Transforms compose outwards: the one on
// the child happens to the shape before the one above it.
func (s *Style) withTransform(raw string, warn func(string, ...any)) {
	m, ok := parseTransform(raw)
	switch {
	case ok:
		if s.HasTransform {
			s.Transform = s.Transform.Mul(m)
		} else {
			s.Transform, s.HasTransform = m, true
		}
	case strings.TrimSpace(raw) == "":
	default:
		warn("the transform %q is not one this package can read, so the shape stays where it is", raw)
	}
}
