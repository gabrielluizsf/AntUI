package svg

import (
	"strconv"
	"strings"
)

// applyStrokeAttrs reads the stroke shape attributes: how the ends and corners
// are treated, the miter limit, whether the stroke is scaled, and the path
// length the stroke is measured out against.
func (s *Style) applyStrokeAttrs(attrs map[string]string, warn func(string, ...any)) {
	for name, raw := range attrs {
		switch name {
		case "stroke-linecap":
			s.Cap = readLineCap(raw)
		case "stroke-linejoin":
			s.Join = readLineJoin(raw)
		case "stroke-miterlimit":
			if v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64); err == nil && v > 0 {
				s.MiterLimit = v
			}
		case "vector-effect":
			s.withVectorEffect(raw, warn)
		case "pathlength":
			s.pathLength = readPathLength(raw, warn)
		}
	}
}
