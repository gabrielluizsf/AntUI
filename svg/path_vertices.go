package svg

import "github.com/gabrielluizsf/antui/canvas"

// pathVertices is every vertex of the path, subpath by subpath, in the order
// the path draws them, each with the direction the path reaches it and leaves
// it by.
//
// A subpath that closes travels back to its first point, so that point is both
// where the path leaves and where it arrives: it is written once with both
// directions, and `orient="auto"` turns its marker to the bisector between them
// rather than to one of the two sides of the corner. A subpath written after a
// close starts at the point the close came back to and is the same shared
// vertex seen from the other side, so it carries the direction the close
// arrived by and leaves along the closed subpath's own first segment — which is
// only known from the subpath before it.
func pathVertices(path *canvas.Path) []markerVertex {
	pts, closed := path.Points()
	var vs []markerVertex
	for j, sub := range pts {
		if len(sub) == 0 {
			continue
		}
		// carried is the closed subpath before this one when this subpath is
		// where that close came back to; a subpath after a move has nothing
		// arriving at its first point, because the path does not travel from
		// the end of one subpath to the start of the next.
		var carried []canvas.Point
		if j > 0 && closed[j-1] && len(pts[j-1]) > 1 {
			carried = pts[j-1]
		}
		for k, q := range sub {
			v := markerVertex{pt: q}
			switch {
			case k > 0:
				v.hasIn, v.inX, v.inY = true, q.X-sub[k-1].X, q.Y-sub[k-1].Y
			case closed[j] && len(sub) > 1:
				last := sub[len(sub)-1]
				v.hasIn, v.inX, v.inY = true, q.X-last.X, q.Y-last.Y
			case carried != nil:
				last := carried[len(carried)-1]
				v.hasIn, v.inX, v.inY = true, q.X-last.X, q.Y-last.Y
			}
			switch {
			case k+1 < len(sub):
				v.hasOut, v.outX, v.outY = true, sub[k+1].X-q.X, sub[k+1].Y-q.Y
			case closed[j] && len(sub) > 1:
				v.hasOut, v.outX, v.outY = true, sub[0].X-q.X, sub[0].Y-q.Y
			case carried != nil:
				v.hasOut, v.outX, v.outY = true, carried[1].X-q.X, carried[1].Y-q.Y
			}
			vs = append(vs, v)
		}
	}
	return vs
}
