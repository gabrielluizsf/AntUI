package svg

import (
	"strconv"
	"strings"
)

// markerLength is one length a marker is written with, read against the
// reference a percentage of it is a fraction of — which differs from attribute
// to attribute and is the caller's to know. That reference is the drawing's
// viewport wherever the length names the room, as it is for a length of the
// drawing's own; SVG 2 leaves it open for these attributes and browsers
// disagree on it, so the viewport is the one this counts. What this cannot read
// answers false rather than a default, so the caller can say so out loud: a
// length quietly taken as the default is one an author never learns was
// ignored.
func markerLength(raw string, ref float64) (float64, bool) {
	if body, ok := strings.CutSuffix(raw, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(body), 64)
		if err != nil {
			return 0, false
		}
		return v / 100 * ref, true
	}
	return parseLength(raw)
}
