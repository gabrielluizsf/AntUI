package css

import (
	"strings"
)

// parseTextOverflow reads the text-overflow keyword. The two-value form
// (clip ellipsis) is accepted with the last keyword winning.
func parseTextOverflow(raw string) (uint8, bool) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return 0, false
	}
	v := TextOverflowClip
	for _, p := range parts {
		switch p {
		case "clip":
			v = TextOverflowClip
		case "ellipsis":
			v = TextOverflowEllipsis
		default:
			return 0, false
		}
	}
	return v, true
}
