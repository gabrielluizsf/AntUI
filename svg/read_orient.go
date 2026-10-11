package svg

import "strings"

// readOrient is how a marker is turned: an angle, or the keyword that turns it
// to the direction of the path at the vertex it is drawn on. An angle without a
// unit is degrees, as everywhere else in SVG. What it does not say is the
// default of no turn at all, and what it says but this package cannot read
// keeps that default with one warning, because a marker pointing the wrong way
// is one an author can see.
func readOrient(raw string, warn func(string, ...any)) (deg float64, auto, startReverse bool) {
	switch s := strings.ToLower(strings.TrimSpace(raw)); s {
	case "":
		return 0, false, false
	case "auto":
		return 0, true, false
	case "auto-start-reverse":
		return 0, true, true
	}
	v, err := parseAngle(raw)
	if err != nil {
		warn("the marker orient %q is not an angle, so the marker is drawn the way it was written", raw)
		return 0, false, false
	}
	return v, false, false
}
