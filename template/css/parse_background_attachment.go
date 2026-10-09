package css

import (
	"strings"
)

// parseBackgroundAttachment reads a comma-separated scroll/fixed/local list.
func parseBackgroundAttachment(raw string) ([]uint8, bool) {
	var out []uint8
	for _, part := range splitFields(raw, ',') {
		switch strings.TrimSpace(part) {
		case "scroll":
			out = append(out, BackAttachScroll)
		case "fixed":
			out = append(out, BackAttachFixed)
		case "local":
			out = append(out, BackAttachLocal)
		default:
			return nil, false
		}
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}
