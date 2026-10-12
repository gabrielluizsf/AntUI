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
	w, h, onto, ok := markerRoom(def, st)
	if !ok {
		return nil
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
	return &markerPlace{def: def, at: at, room: markerClip(onto, w, h, def.visible), prop: prop}
}
