package canvas

import "math"

// Point is one vertex of a path, in whatever units the path was built in —
// user units, a viewBox's coordinates, a reference-pixel grid. A path is
// geometry; only the canvas it is drawn on knows what a unit costs.
type Point struct {
	X, Y float64
}

// FillRule decides which parts of a self-overlapping path are inside it.
type FillRule uint8

// The two rules SVG and CSS ask for. NonZero is the default: a shape counts how
// many times it winds around a point and is inside wherever the count is not
// zero, so two nested rings drawn the same way are a ring and a hole. EvenOdd
// ignores the direction and counts the crossings, so any two nested rings leave
// a hole between them, however they were wound.
const (
	FillNonZero FillRule = iota
	FillEvenOdd
)

// Path is a sequence of subpaths. Each subpath is a run of points that may or
// may not be closed, and a fill treats every subpath as a boundary while a
// stroke walks each one as a line of its own. Curves are flattened into points
// when they are added, so a Path is a polygon set by the time anything reads it.
type Path struct {
	subs []subpath
}

// subpath is one run of points and whether the last one joins the first.
type subpath struct {
	pts    []Point
	closed bool
}

// NewPath makes an empty path.
func NewPath() *Path { return &Path{} }

// Empty reports whether the path holds no subpath, or only subpaths that never
// got a point. A path with nothing in it draws nothing at all, and saying so
// early saves walking it.
func (p *Path) Empty() bool {
	for _, s := range p.subs {
		if len(s.pts) > 0 {
			return false
		}
	}
	return true
}

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

// MoveTo starts a new subpath at x, y. A move with no line behind it still
// opens a subpath, so a path of only moves fills nothing but strokes dots.
func (p *Path) MoveTo(x, y float64) {
	p.subs = append(p.subs, subpath{pts: []Point{{x, y}}})
}

// LineTo adds a point to the subpath being drawn, starting one at the point
// itself when the path is empty — the first command of a path needs no move to
// have come before it.
func (p *Path) LineTo(x, y float64) {
	if len(p.subs) == 0 {
		p.MoveTo(x, y)
		return
	}
	s := &p.subs[len(p.subs)-1]
	s.pts = append(s.pts, Point{x, y})
}

// Close joins the last subpath back to its first point and opens a new subpath
// at that same point, so a line written after the close carries on from where
// the path closed — which is what a moveto-less "Z L" pair means.
func (p *Path) Close() {
	s := p.last()
	if s == nil || len(s.pts) == 0 {
		return
	}
	s.closed = true
	p.subs = append(p.subs, subpath{pts: []Point{s.pts[0]}})
}

// last returns the subpath being drawn, or nil when the path is empty.
func (p *Path) last() *subpath {
	if len(p.subs) == 0 {
		return nil
	}
	return &p.subs[len(p.subs)-1]
}

// end is the point the path is currently at.
func (p *Path) end() Point {
	s := p.last()
	if s == nil || len(s.pts) == 0 {
		return Point{}
	}
	return s.pts[len(s.pts)-1]
}

// Transform maps every point of the path through m, in place. Flattening
// happens before it, so a curve rotated into place is rotated as the polygon
// it already became — which is the flattening this rasterizer is honest about.
func (p *Path) Transform(m Matrix) {
	for _, s := range p.subs {
		for i := range s.pts {
			s.pts[i].X, s.pts[i].Y = m.Map(s.pts[i].X, s.pts[i].Y)
		}
	}
}

// Bounds is the tightest box holding every point, and reports whether the path
// has any point at all to bound.
func (p *Path) Bounds() (minX, minY, maxX, maxY float64, ok bool) {
	for _, s := range p.subs {
		for _, pt := range s.pts {
			if !ok {
				minX, minY, maxX, maxY, ok = pt.X, pt.Y, pt.X, pt.Y, true
				continue
			}
			minX = min(minX, pt.X)
			minY = min(minY, pt.Y)
			maxX = max(maxX, pt.X)
			maxY = max(maxY, pt.Y)
		}
	}
	return minX, minY, maxX, maxY, ok
}

// Length is how far it is along the path: every segment of every subpath, with
// a closed one counted back to the point it started from. Curves were flattened
// into points as they were added, so this measures the polygon the path has
// become — the same walk a stroke takes along it and the same walk a dash is
// cut into. A path with nothing in it, or only points that lead nowhere, is
// nothing long.
func (p *Path) Length() float64 {
	if p == nil {
		return 0
	}
	pts, closed := p.Points()
	total := 0.0
	for i, sub := range pts {
		for j := 1; j < len(sub); j++ {
			total += math.Hypot(sub[j].X-sub[j-1].X, sub[j].Y-sub[j-1].Y)
		}
		if closed[i] && len(sub) > 1 {
			from, to := sub[len(sub)-1], sub[0]
			total += math.Hypot(to.X-from.X, to.Y-from.Y)
		}
	}
	return total
}
