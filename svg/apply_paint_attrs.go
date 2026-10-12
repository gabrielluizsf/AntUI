package svg

import "strings"

// applyPaintAttrs reads the paint attributes: the fill and stroke paints and
// their opacities, the clip-path, clip-rule, mask and filter references, and
// the markers.
func (s *Style) applyPaintAttrs(attrs map[string]string, warn func(string, ...any)) {
	for name, raw := range attrs {
		switch name {
		case "fill":
			s.withFill(raw)
		case "fill-rule":
			s.FillRule = readFillRule(raw)
		case "clip-path":
			s.withClipPath(raw, warn)
		case "clip-rule":
			s.ClipRule = readFillRule(raw)
		case "mask":
			s.withMask(raw, warn)
		case "filter":
			s.filters = readFilters(raw, warn)
		case "marker-start", "marker-mid", "marker-end":
			s.withMarker(name, raw, warn)
		case "fill-opacity":
			if v, ok := parseAlpha(raw); ok {
				s.FillOpacity = v
			}
		case "stroke":
			s.withStroke(raw)
		case "stroke-width":
			if v, ok := parseLength(raw); ok {
				s.Width = v
			} else if strings.TrimSpace(raw) != "" {
				warn("the stroke width %q is not a length, so the stroke keeps the width it had", raw)
			}
		case "stroke-opacity":
			if v, ok := parseAlpha(raw); ok {
				s.StrokeOpacity = v
			}
		case "opacity":
			if v, ok := parseAlpha(raw); ok {
				s.Opacity = v
			}
		}
	}
}
