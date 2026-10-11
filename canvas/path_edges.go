package canvas

import "slices"

// pathEdges turns every subpath into the edges a line can cross: each is a run
// of points joined back up to its first when closed, and left out when it is a
// single point, which has no area to cross.
func pathEdges(path *Path, y0, y1 int) []edge {
	pts, closed := path.Points()
	edges := make([]edge, 0, len(pts))
	for i, sub := range pts {
		if len(sub) < 2 {
			continue
		}
		if closed[i] {
			sub = slices.Concat(sub, sub[:1])
		}
		for j := 1; j < len(sub); j++ {
			if e, ok := makeEdge(sub[j-1], sub[j], y0, y1); ok {
				edges = append(edges, e)
			}
		}
	}
	return edges
}
