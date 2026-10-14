package svg

import "strings"

// withTransform reads a `transform`. Transforms compose outwards: the one on
// the child happens to the shape before the one above it. The shift the matrix
// writes is kept beside the composition because fixed-position spends it as the
// origin of the fixed user coordinate system.
func (s *Style) withTransform(raw string, warn func(string, ...any)) {
	m, ok := parseTransform(raw)
	switch {
	case ok:
		s.VectorEffect.OriginX, s.VectorEffect.OriginY = m.E, m.F
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
