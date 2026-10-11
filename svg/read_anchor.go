package svg

import "strings"

// readAnchor reads where along the pen a piece of writing hangs: at it, centred
// on it, or ending at it. The third of those is what makes a label sit to the
// left of the point it names rather than to the right.
func readAnchor(raw string) (TextAnchor, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "start":
		return AnchorStart, true
	case "middle":
		return AnchorMiddle, true
	case "end":
		return AnchorEnd, true
	}
	return AnchorStart, false
}
