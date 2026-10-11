package svg

import "strings"

// readMarker reads one `<marker>`'s own sayings — the room it asks for, the
// point of its picture that lands on the vertex, and which way it is turned —
// into the place its id already put. What it does not say is the SVG default,
// which is a room three stroke-widths square with the reference at its origin
// and no turn on it.
//
// The style the picture is built from starts at the drawing's own rather than
// at the style of the place the marker stands, which is what [Image.readClipPath]
// says about a clipPath too. The spec has properties inherit into a marker from
// its ancestors and not from the element that names it, and the second half of
// that is the half that matters here: a marker is drawn at a vertex far from
// wherever it was written, in colours of its own, and taking the colours of the
// shape that happens to name it would make the same marker come out differently
// at every end of every line. The first half is what this leaves out — a marker
// under a group that painted everything in it does not take that paint — and a
// marker whose picture wants a colour says so itself, as they are written.
func (img *Image) readMarker(def *markerDef, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	if raw := strings.TrimSpace(e.attr("viewBox")); raw != "" {
		if v, ok := parseViewBox(raw); ok {
			def.view, def.hasView = v, true
		} else {
			warn("the marker viewBox %q is not four numbers, so its picture is drawn where it was written", raw)
		}
	}
	def.refX = markerLength(e, "refX", 0, warn)
	def.refY = markerLength(e, "refY", 0, warn)
	def.w = markerLength(e, "markerWidth", 3, warn)
	def.h = markerLength(e, "markerHeight", 3, warn)
	def.strokeUnits = !strings.EqualFold(strings.TrimSpace(e.attr("markerUnits")), "userspaceonuse")
	def.deg, def.auto, def.startReverse = readOrient(e.attr("orient"), warn)
	def.visible = overflowIn(e)
	if def.w <= 0 || def.h <= 0 {
		warn("the marker is %g by %g, which has no room for anything in it, so nothing is drawn where it is named", def.w, def.h)
	}
}
