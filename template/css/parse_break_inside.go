package css

import (
	"strings"
)

// parseBreakInside reads whether a box may be cut in the middle. The
// page-break-inside name the property had before break-inside replaced it means
// the same thing here, so both reach the same field. avoid-column, avoid-page
// and avoid-region all say avoid: the engine has one kind of fragment, so it
// has one answer to all of them.
func parseBreakInside(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "auto":
		return BreakAuto, true
	case "avoid", "avoid-column", "avoid-page", "avoid-region":
		return BreakAvoid, true
	}
	return 0, false
}
