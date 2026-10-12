package svg

import "strings"

// applyTextAttrs reads the writing attributes: the size, anchor, family,
// weight and letter spacing, the whitespace mode, whether the element is
// visible at all, and the transform.
func (s *Style) applyTextAttrs(attrs map[string]string, warn func(string, ...any)) {
	for name, raw := range attrs {
		switch name {
		case "font-size":
			s.FontSize = readFontSize(s.FontSize, raw, warn)
		case "text-anchor":
			if a, ok := readAnchor(raw); ok {
				s.Anchor = a
			} else if strings.TrimSpace(raw) != "" {
				warn("the text-anchor %q is not start, middle or end, so the writing stays where the pen is", raw)
			}
		case "xml:space":
			s.PreserveSpace = readPreserveSpace(raw, s.PreserveSpace)
		case "font-family":
			s.FontFamily = strings.TrimSpace(raw)
		case "font-weight":
			s.FontWeight = readFontWeight(s.FontWeight, raw, warn)
		case "letter-spacing":
			s.LetterSpacing = readLetterSpacing(raw, warn)
		case "display", "visibility":
			if strings.EqualFold(strings.TrimSpace(raw), "none") ||
				strings.EqualFold(strings.TrimSpace(raw), "hidden") {
				s.Hidden = true
			}
		case "transform":
			s.withTransform(raw, warn)
		}
	}
}
