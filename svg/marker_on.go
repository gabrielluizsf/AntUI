package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// markerOn is one marker on one vertex of a shape: the room the picture is
// drawn into, the fit of its viewBox, the turn `orient` gave it, and the matrix
// that puts all of that down on the vertex — with the element's own transform
// in it under the same condition the path itself got it, so that a marker lands
// where the vertex landed.
//
// The picture is taken there the way it was written: fitted into the room the
// marker asked for if it has a viewBox, moved back by the reference point,
// turned, and put down on the vertex. At the first vertex a marker written
// `auto-start-reverse` is turned around, which is what lets one marker be put
// at both ends of a line and point out of each of them; everywhere else it is
// the same as `auto`, as the spec says.
func markerOn(def *markerDef, st Style, v markerVertex, atStart bool, prop string) *markerPlace {
	if def == nil {
		return nil
	}
	// The room the picture is drawn into: the size the marker asked for, and a
	// multiple of the stroke-width of this shape where markerUnits says so —
	// the default, which is what makes a marker on a thick line a thick
	// marker. Without a viewBox the room only bounds what is drawn rather than
	// scaling it, since the picture is written in the room's own coordinates
	// already; with one, the fit below is what the two numbers are used for.
	scale := 1.0
	if def.strokeUnits {
		scale = st.Width
	}
	w, h := def.w*scale, def.h*scale
	if !(w > 0) || !(h > 0) {
		return nil
	}
	// The picture onto the room, which is the fit a viewBox gets everywhere it
	// is fitted to something — the canvas, the box a `<use>` asked a `<symbol>`
	// for — kept proportions and all.
	onto := canvas.Identity()
	if def.hasView {
		onto = fitTransform(def.view, w, h)
	}
	rx, ry := onto.Map(def.refX, def.refY)
	angle := def.deg * math.Pi / 180
	if def.auto {
		angle = v.angle()
		if atStart && def.startReverse {
			angle += math.Pi
		}
	}
	at := canvas.Translate(v.pt.X, v.pt.Y).
		Mul(canvas.Rotate(angle)).
		Mul(canvas.Translate(-rx, -ry)).
		Mul(onto)
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		at = st.Transform.Mul(at)
	}
	// The room as a cut, in the coordinates the picture is written in: its
	// rectangle taken backwards through the fit of the viewBox, so that where
	// the room's edges land on the picture follow where the picture's own
	// writing lands. A fit with no way back collapses the picture to nothing
	// anyway, so it gets no cut rather than the wrong one.
	var room *clipRegion
	if !def.visible {
		if inv, ok := onto.Inverse(); ok {
			p := canvas.NewPath()
			p.AddRect(0, 0, w, h)
			p.Transform(inv)
			room = &clipRegion{shapes: []canvas.MaskShape{{Path: p, Rule: canvas.FillNonZero}}}
		}
	}
	return &markerPlace{def: def, at: at, room: room, prop: prop}
}
