package css

import "strings"

// parseGridTemplate reads a whole grid-template shorthand.
func parseGridTemplate(raw string, ctx Units) (GridTemplate, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "none") {
		return nil, true
	}
	return parseGridTemplateTracks(raw, ctx, 0)
}

// parseGridTemplateTracks reads one track list — the value itself, or the body
// of a repeat() — into the segments it stands for. A "[…]" group names the
// line to its left: names seen before any track wait for it, and names seen
// after one close the line that track ends.
