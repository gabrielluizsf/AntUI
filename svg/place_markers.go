package svg

import "github.com/gabrielluizsf/antui/canvas"

// placeMarkers is where the three markers an element asked for sit on its
// shape, one entry for each marker it draws in the order the vertices come:
// the first vertex for marker-start, every vertex between the first and the
// last for marker-mid, and the last one for marker-end. A path of a single
// vertex is both its first and its last, so it takes a marker-start and a
// marker-end in the same place and no marker-mid at all.
//
// It answers nothing where the shape has no path, where no marker was named —
// which is all of the shapes in most of the drawings — or where the room a
// marker asked for is no room at all, which draws nothing where it is named.
func placeMarkers(path *canvas.Path, st Style) []markerPlace {
	if path == nil || (st.markerStart == nil && st.markerMid == nil && st.markerEnd == nil) {
		return nil
	}
	vs := pathVertices(path)
	if len(vs) == 0 {
		return nil
	}
	var out []markerPlace
	if p := markerOn(st.markerStart, st, vs[0], true, "marker-start"); p != nil {
		out = append(out, *p)
	}
	if st.markerMid != nil {
		for i := 1; i+1 < len(vs); i++ {
			if p := markerOn(st.markerMid, st, vs[i], false, "marker-mid"); p != nil {
				out = append(out, *p)
			}
		}
	}
	if p := markerOn(st.markerEnd, st, vs[len(vs)-1], false, "marker-end"); p != nil {
		out = append(out, *p)
	}
	return out
}
