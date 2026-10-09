package css

import "strings"

// parseGridAutoRepeat reads the "auto-fit" or "auto-fill" keyword that introduces an auto-repeated track list.
func parseGridAutoRepeat(raw string) (fit bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "auto-fill":
		return false, true
	case "auto-fit":
		return true, true
	}
	return false, false
}

// parseGridAutoRepeatTracks reads the track an auto-repeat stands for: one
// track, no line names, and a minimum small enough to count against the space
// the axis has.
