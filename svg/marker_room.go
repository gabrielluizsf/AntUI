package svg

import "github.com/gabrielluizsf/antui/canvas"

// markerRoom is the room a marker's picture is drawn into — the size the marker
// asked for, and a multiple of the stroke-width of this shape where markerUnits
// says so, which is the default and what makes a marker on a thick line a thick
// marker — and the fit of its viewBox onto it. Without a viewBox the room only
// bounds what is drawn rather than scaling it, since the picture is written in
// the room's own coordinates already. It answers false where the room has no
// area to draw into.
func markerRoom(def *markerDef, st Style) (w, h float64, onto canvas.Matrix, ok bool) {
	scale := 1.0
	if def.strokeUnits {
		scale = st.Width
	}
	w, h = def.w*scale, def.h*scale
	if !(w > 0) || !(h > 0) {
		return 0, 0, canvas.Matrix{}, false
	}
	// The picture onto the room, which is the fit a viewBox gets everywhere it
	// is fitted to something — the canvas, the box a `<use>` asked a `<symbol>`
	// for — kept proportions and all.
	onto = canvas.Identity()
	if def.hasView {
		onto = fitTransform(def.view, w, h)
	}
	return w, h, onto, true
}
