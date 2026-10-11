package svg

// resolveMarkers follows the three marker references an element named, which
// only the whole drawing can say anything about, and answers the markers that
// element draws its vertices with: the `<marker>` each id names, or one
// warning and nothing where the drawing does not have it.
//
// Each reference is cleared as it is spent, the same way a fill's is, so that a
// group says what it has to say once and not once for every shape inside it.
// What the references turn into are NOT cleared, which is how this differs from
// a clip, a mask and a filter: marker-start, marker-mid and marker-end are
// inherited, so a `<g marker-end>` is what every shape under it draws a marker
// with at its last vertex, and the children take the marker itself rather than
// the id to look up again. An element that writes a marker of its own —
// including `none`, which turns it off — takes a new one over in [Style.with],
// and that is where the old one is dropped.
func (img *Image) resolveMarkers(st *Style, warn func(string, ...any)) {
	for _, r := range []struct {
		prop string
		ref  *string
		def  **markerDef
	}{
		{"marker-start", &st.markerStartRef, &st.markerStart},
		{"marker-mid", &st.markerMidRef, &st.markerMid},
		{"marker-end", &st.markerEndRef, &st.markerEnd},
	} {
		if *r.ref == "" {
			continue
		}
		id := *r.ref
		*r.ref = ""
		def := img.markers[id]
		if def == nil {
			warn("the %s names #%s, which the drawing does not have, so no marker is drawn where it points", r.prop, id)
			continue
		}
		*r.def = def
	}
}
