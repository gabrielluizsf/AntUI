package svg

import "strings"

// withMarker reads one of `marker-start`, `marker-mid` and `marker-end`: the
// markers drawn at the vertices of this element's shape. Unlike the clip, the
// mask and the filter they ARE inherited, so the reference is kept for the
// whole drawing to follow and what was inherited is dropped as this element
// writes one of its own — `none` among them, which is how a shape turns off
// the marker its group asked for. Each drops only the one it writes: a shape
// that says `marker-start="none"` still carries the marker-mid the group above
// it asked for. See [Image.resolveMarkers].
func (s *Style) withMarker(name, raw string, warn func(string, ...any)) {
	ref, def := &s.markerStartRef, &s.markerStart
	if name == "marker-mid" {
		ref, def = &s.markerMidRef, &s.markerMid
	} else if name == "marker-end" {
		ref, def = &s.markerEndRef, &s.markerEnd
	}
	*ref, *def = "", nil
	switch id, ok := paintRef(raw); {
	case strings.EqualFold(strings.TrimSpace(raw), "none"):
	case ok:
		*ref = id
	default:
		warn("the %s %q is not a reference to a <marker> this package can follow, so no marker is drawn at the vertices it names", name, raw)
	}
}
