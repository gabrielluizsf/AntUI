package svg

import "strings"

// markerLength is one length a marker is written with: what it does not say is
// the default the spec gives it, and what it says but this package cannot read
// is said out loud rather than quietly taken as the default. A percentage is
// left out of that — it is a fraction of the room the marker is drawn into,
// which is not known until the shape that names it is — so it comes through
// here as something that is not a length and is said.
func markerLength(e *element, name string, def float64, warn func(string, ...any)) float64 {
	if !e.hasAttr(name) {
		return def
	}
	raw := strings.TrimSpace(e.attr(name))
	v, ok := parseLength(raw)
	if !ok {
		warn("the marker %s %q is not a length, so it is taken as %g", name, raw, def)
		return def
	}
	return v
}
