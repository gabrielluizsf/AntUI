package svg

import "strings"

// readMarker reads one `<marker>`'s own sayings — the room it asks for, the
// point of its picture that lands on the vertex, which way it is turned, and
// how its viewBox is fitted — into the place its id already put. What it does
// not say is the SVG default, which is a room three stroke-widths square with
// the reference at its origin and no turn on it.
//
// The style the picture is built from is the one that reaches the `<marker>`
// from the elements containing it, which is read in [Image.readMarkers] as the
// walk passes and applied where the picture is built. It is not the style of
// the element that names the marker: a marker is drawn at a vertex far from
// wherever it was written, and taking the colours of the shape that happens to
// name it would make the same marker come out differently at every end of every
// line. A picture that wants those colours asks for them by name, as
// `context-fill` and `context-stroke`, rather than inheriting them.
func (img *Image) readMarker(def *markerDef, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	if raw := strings.TrimSpace(e.attr("viewBox")); raw != "" {
		if v, ok := parseViewBox(raw); ok {
			def.view, def.hasView = v, true
		} else {
			warn("the marker viewBox %q is not four numbers, so its picture is drawn where it was written", raw)
		}
	}
	// A percentage counts against the drawing's viewport, as a length of the
	// drawing's own does — see [markerLength] for why that is the reference.
	length := func(name string, ref, dflt float64) float64 {
		if !e.hasAttr(name) {
			return dflt
		}
		raw := strings.TrimSpace(e.attr(name))
		v, ok := markerLength(raw, ref)
		if !ok {
			warn("the marker %s %q is not a length, so it is taken as %g", name, raw, dflt)
			return dflt
		}
		return v
	}
	def.w = length("markerWidth", img.ViewBox[2], 3)
	def.h = length("markerHeight", img.ViewBox[3], 3)
	markerRefs(def, e, warn)
	def.strokeUnits = !strings.EqualFold(strings.TrimSpace(e.attr("markerUnits")), "userspaceonuse")
	def.deg, def.auto, def.startReverse = readOrient(e.attr("orient"), warn)
	def.visible = overflowIn(e)
	// preserveAspectRatio decides the fit of a viewBox, so without one it is
	// passed over in silence, as the spec has it, and the default stands.
	def.fit, def.align = PictureFitMeet, AlignXMidYMid
	if def.hasView {
		def.fit, def.align = readPreserveAspect(e.attr("preserveAspectRatio"), warn)
	}
	if def.w <= 0 || def.h <= 0 {
		warn("the marker is %g by %g, which has no room for anything in it, so nothing is drawn where it is named", def.w, def.h)
	}
}
