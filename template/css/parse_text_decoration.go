package css

import (
	"strings"
)

// parseTextDecoration reads the text-decoration shorthand or its -line
// longhand: none clears the flags, otherwise underline, overline and
// line-through accumulate. Style and colour words the engine ignores are
// skipped, not rejected.
func parseTextDecoration(raw string) (uint8, bool) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return 0, false
	}
	var bits uint8
	for _, p := range parts {
		switch p {
		case "none":
			return 0, true
		case "underline":
			bits |= TextDecorationUnderline
		case "overline":
			bits |= TextDecorationOverline
		case "line-through":
			bits |= TextDecorationLineThrough
		default:
			// Thickness, style and colour words are not modelled; skip them.
		}
	}
	return bits, true
}
