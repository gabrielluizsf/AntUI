package svg

import "strings"

// markerRefs reads the point of the picture that lands on the vertex: refX and
// refY, a percentage of the viewBox where the marker has one — the spec's own
// words — and of the room where it has not, which is the same thing, since
// without a viewBox the picture is written in the room's own coordinates.
func markerRefs(def *markerDef, e *element, warn func(string, ...any)) {
	w, h := def.w, def.h
	if def.hasView {
		w, h = def.view[2], def.view[3]
	}
	ref := func(name string, base float64, keys map[string]float64) float64 {
		if !e.hasAttr(name) {
			return 0
		}
		raw := strings.TrimSpace(e.attr(name))
		if f, ok := keys[strings.ToLower(raw)]; ok {
			return f * base
		}
		v, ok := markerLength(raw, base)
		if !ok {
			warn("the marker %s %q is not a length, so it is taken as 0", name, raw)
		}
		return v
	}
	def.refX = ref("refX", w, markerRefXKeys)
	def.refY = ref("refY", h, markerRefYKeys)
}
