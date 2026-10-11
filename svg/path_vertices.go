package svg

import "github.com/gabrielluizsf/antui/canvas"

// pathVertices is every vertex of the path, subpath by subpath, in the order
// the path draws them. The point a close travelled back to is a vertex like any
// other, which is what makes the closing vertex of a rectangle its last one and
// the same place as its first — the way SVG says the last vertex of a path that
// ends closed is its initial vertex again.
func pathVertices(path *canvas.Path) []markerVertex {
	pts, closed := path.Points()
	var vs []markerVertex
	for j, sub := range pts {
		if len(sub) == 0 {
			continue
		}
		// A subpath of one point standing after a closed one is the point that
		// close travelled back to — it is what [canvas.Path.Close] leaves
		// behind — so the path arrives at it along the segment that closed the
		// subpath before it.
		remnant := len(sub) == 1 && j > 0 && closed[j-1] && len(pts[j-1]) > 0
		for k, q := range sub {
			v := markerVertex{pt: q}
			switch {
			case k > 0:
				v.hasIn, v.inX, v.inY = true, q.X-sub[k-1].X, q.Y-sub[k-1].Y
			case remnant:
				prev := pts[j-1]
				v.hasIn, v.inX, v.inY = true, q.X-prev[len(prev)-1].X, q.Y-prev[len(prev)-1].Y
			}
			switch {
			case k+1 < len(sub):
				v.hasOut, v.outX, v.outY = true, sub[k+1].X-q.X, sub[k+1].Y-q.Y
			case closed[j] && len(sub) > 1:
				v.hasOut, v.outX, v.outY = true, sub[0].X-q.X, sub[0].Y-q.Y
			}
			vs = append(vs, v)
		}
	}
	return vs
}
