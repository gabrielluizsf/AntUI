package svg

import "github.com/gabrielluizsf/antui/canvas"

// markerRoom is the room a marker's picture is drawn into — the size the marker
// asked for, and a multiple of the stroke-width of this shape where markerUnits
// says so, which is the default and what makes a marker on a thick line a thick
// marker — and the fit of its picture onto it. A viewBox is fitted into the room
// the way the marker's own preserveAspectRatio asks; without one the picture is
// written in the marker's own units, which that same multiple is the size of,
// so it is drawn at the multiple and only the room's edge cuts what runs out.
// It answers false where the room has no area to draw into.
func markerRoom(def *markerDef, st Style) (w, h float64, onto canvas.Matrix, ok bool) {
	scale := 1.0
	if def.strokeUnits {
		scale = st.Width
	}
	w, h = def.w*scale, def.h*scale
	if !(w > 0) || !(h > 0) {
		return 0, 0, canvas.Matrix{}, false
	}
	onto = canvas.Scale(scale, scale)
	if def.hasView {
		onto = fitViewBox(def.view, &Picture{
			Box:   [4]float64{0, 0, w, h},
			Fit:   def.fit,
			Align: def.align,
		})
	}
	return w, h, onto, true
}
