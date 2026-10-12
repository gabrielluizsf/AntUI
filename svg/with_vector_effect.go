package svg

import "strings"

// withVectorEffect reads `vector-effect`: whether the stroke is taken through
// everything that scaled the shape — which is what SVG does by default — or
// drawn at the width it says, in the pixels the drawing comes out at whatever
// scaled it. Any other value is a vector effect this package does not do, and
// saying so beats quietly painting a stroke nobody asked for.
func (s *Style) withVectorEffect(raw string, warn func(string, ...any)) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "none":
		s.NonScalingStroke = false
	case "non-scaling-stroke":
		s.NonScalingStroke = true
	default:
		s.NonScalingStroke = false
		if strings.TrimSpace(raw) != "" {
			warn("the vector-effect %q is not none or non-scaling-stroke, so the stroke is scaled the way SVG scales it", raw)
		}
	}
}
