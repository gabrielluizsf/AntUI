package canvas

// Points hands out the flattened points of every subpath, and whether that
// subpath is closed. The slice is the path's own storage: a caller that keeps
// it must not build on the path afterwards.
func (p *Path) Points() (pts [][]Point, closed []bool) {
	pts = make([][]Point, len(p.subs))
	closed = make([]bool, len(p.subs))
	for i, s := range p.subs {
		pts[i] = s.pts
		closed[i] = s.closed
	}
	return pts, closed
}
