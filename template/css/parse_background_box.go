package css

import (
	"strings"
)

// parseBackgroundBox reads a comma-separated background-clip or
// background-origin list of border/padding/content box keywords.
func parseBackgroundBox(raw string) ([]uint8, bool) {
	var out []uint8
	for _, part := range splitFields(raw, ',') {
		switch strings.TrimSpace(part) {
		case "border-box":
			out = append(out, BackBorder)
		case "padding-box":
			out = append(out, BackPadding)
		case "content-box":
			out = append(out, BackContent)
		default:
			return nil, false
		}
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}
