package svg

import "github.com/gabrielluizsf/antui/canvas"

// applyFilterDef runs a `<filter>` element's primitives over the picture it
// was pointed at, one after another: each takes the pictures its inputs name
// — the picture so far when it names none, the picture the element drew for
// SourceGraphic, and for anything an earlier primitive left behind — and what
// it leaves under its own name is what the rest may take. What comes out is a
// picture of the whole layer, the same as one CSS function's own picture.
func applyFilterDef(dst *canvas.Canvas, def *filterDef, s scale2, current canvas.Color) *canvas.Canvas {
	if dst == nil || len(def.prims) == 0 {
		return dst
	}
	src := dst
	named := map[string]*canvas.Canvas{}
	var withAlpha *canvas.Canvas
	prev := dst
	for _, p := range def.prims {
		in := filterInput(p.in, prev, src, &withAlpha, named)
		in2 := filterInput(p.in2, src, src, &withAlpha, named)
		out := primPicture(p, in, in2, src, s, current)
		if out == nil {
			continue
		}
		prev = out
		if p.result != "" {
			named[p.result] = out
		}
	}
	return prev
}
