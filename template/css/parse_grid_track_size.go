package css

import (
	"math"
	"strconv"
	"strings"
)

// parseGridTrackSize reads one track length, optionally clamped to a flexible weight.
func parseGridTrackSize(raw string, ctx Units, flexible bool) (GridTrackSize, bool) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	switch raw {
	case "auto":
		return GridTrackSize{Kind: GridTrackAuto}, true
	case "min-content":
		return GridTrackSize{Kind: GridTrackMinContent}, true
	case "max-content":
		return GridTrackSize{Kind: GridTrackMaxContent}, true
	}
	if name, inner, ok := gridFunction(raw); ok && name == "fit-content" {
		limit, err := parseLengthAt(inner, ctx)
		if err != nil || limit.Auto() || limit.None() || limit.value < 0 {
			return GridTrackSize{}, false
		}
		return GridTrackSize{Kind: GridTrackFitContent, Length: limit}, true
	}
	if flexible && strings.HasSuffix(raw, "fr") {
		v, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(raw, "fr")), 64)
		if err != nil || v <= 0 || math.IsNaN(v) || math.IsInf(v, 0) {
			return GridTrackSize{}, false
		}
		return GridTrackSize{Kind: GridTrackFlexible, Fr: v}, true
	}
	l, err := parseLengthAt(raw, ctx)
	if err != nil || l.Auto() || l.None() || l.value < 0 {
		return GridTrackSize{}, false
	}
	return GridTrackSize{Kind: GridTrackLength, Length: l}, true
}
