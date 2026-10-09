package css

import (
	"strings"
)

// parseFloat reads the side a box is taken out of the flow towards.
func parseFloat(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "none":
		return FloatNone, true
	case "left":
		return FloatLeft, true
	case "right":
		return FloatRight, true
	case "inline-start":
		return FloatInlineStart, true
	case "inline-end":
		return FloatInlineEnd, true
	}
	return 0, false
}
