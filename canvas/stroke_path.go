package canvas

import (
	"math"
	"slices"
)

// LineCap is what a stroke leaves at the end of an open subpath.
type LineCap uint8

// A butt cap stops exactly at the end point, a round one reaches half the
// width past it, and a square one reaches that same distance to the side as
// well — which is what keeps a dashed line the same length however it ends.
const (
	CapButt LineCap = iota
	CapRound
	CapSquare
)

// LineJoin is how a stroke turns a corner between two segments.
type LineJoin uint8

// A miter join cuts the outside of the corner off square, a round one rounds it
// off, and a bevel one chops it flat. A miter longer than its limit falls back
// to the bevel, because a corner drawn out to a spike many times its own width
// is not what "sharp" was meant to say.
const (
	JoinMiter LineJoin = iota
	JoinRound
	JoinBevel
)

// DefaultMiterLimit is the corner SVG gives up on once a join's miter grows past
// this many stroke widths. Four keeps the spike of a right angle and drops it
// on anything sharper.
const DefaultMiterLimit = 4

// StrokeStyle is everything about a stroke that is not its colour: how wide, how
// the ends and the corners are treated, and how far past that limit a miter may
// grow before it is beveled instead.
type StrokeStyle struct {
	Width      float64
	Cap        LineCap
	Join       LineJoin
	MiterLimit float64
}

// StrokePath draws the outline of a path: a band of the given width centred on
// the line, looking the same width whichever way the line runs because every
// piece of the band is filled the same way. A width of zero or less draws
// nothing, which is what SVG says and the only sensible reading of it.
func (cv *Canvas) StrokePath(path *Path, c Color, style StrokeStyle) {
	cv.FillPath(StrokeOutline(path, style), c, FillNonZero)
}

// StrokeOutline is the shape a stroke covers, ready to be filled: what a
// shadow, a glow or a second stroke laid over the first one needs. Every piece
// of it is wound the same way round, because the fill is nonzero and two pieces
// that overlapped the wrong way over would cancel each other out and leave a
// hole where the join should be. It is nil when there is nothing to draw.
func StrokeOutline(path *Path, style StrokeStyle) *Path {
	if path == nil || style.Width <= 0 {
		return nil
	}
	limit := style.MiterLimit
	if limit <= 0 {
		limit = DefaultMiterLimit
	}
	pts, closed := path.Points()
	outline := NewPath()
	half := style.Width / 2
	for i, sub := range pts {
		// Compacted into a copy, because the path belongs to the caller and
		// squeezing runs of the same point out of it in place would show.
		sub = slices.Compact(slices.Clone(sub))
		switch {
		case len(sub) < 2:
			// A lone point is a line with no length, and a cap is all that
			// decides what a stroke of one covers: a round or square cap still
			// reaches past the end, so the point is a dot, and a butt cap stops
			// at it and draws nothing.
			if dot := zeroLengthCap(sub, half, style.Cap); len(dot) > 0 {
				outline.AddPolyline(wound(dot), true)
			}
		case closed[i]:
			for _, piece := range ringOutline(sub, half, style, limit) {
				outline.AddPolyline(wound(piece), true)
			}
		default:
			for _, piece := range openOutline(sub, half, style, limit) {
				outline.AddPolyline(wound(piece), true)
			}
		}
	}
	if outline.Empty() {
		return nil
	}
	return outline
}

// zeroLengthCap is what a subpath that came down to a single point covers. It is
// the cap drawn at that point with no direction to reach along, so a round one is
// a disc and a square one a box, while a butt cap leaves nothing to fill.
func zeroLengthCap(sub []Point, half float64, cap LineCap) []Point {
	if len(sub) != 1 {
		return nil
	}
	switch cap {
	case CapRound:
		return capDisc(sub[0], half)
	case CapSquare:
		return capOutline(sub[0], Point{1, 0}, half, cap)
	}
	return nil
}

// ringOutline is the band around a closed loop: a rectangle along every segment
// and the wedge the two rectangles meeting at a corner leave between them. The
// pieces are kept apart rather than stitched into one loop, because a stroke is
// filled as a union: walking from one rectangle to the next would draw seams
// across the band, and those seams cancel each other out under the nonzero rule.
func ringOutline(sub []Point, half float64, style StrokeStyle, limit float64) [][]Point {
	var out [][]Point
	for i := range sub {
		a, b := sub[i], sub[(i+1)%len(sub)]
		if a == b {
			continue
		}
		out = append(out, segmentOutline(a, b, half))
		// The corner is at the far end of this segment, which is where the
		// next one starts, so the join needs the two points around it.
		if wedge := joinOutline(a, b, sub[(i+2)%len(sub)], half, style, limit); len(wedge) > 0 {
			out = append(out, wedge)
		}
	}
	return out
}

// openOutline is the band along an open run: the same rectangles and joins,
// with a cap at either end so the line stops where it was told to. It is a list
// of pieces rather than one walk, because the fill is a union and a piece does
// not have to be connected to the next one to be filled.
func openOutline(sub []Point, half float64, style StrokeStyle, limit float64) [][]Point {
	var out [][]Point
	out = append(out, segmentOutline(sub[0], sub[1], half))
	for i := 1; i+1 < len(sub); i++ {
		out = append(out, segmentOutline(sub[i], sub[i+1], half))
		out = append(out, joinOutline(sub[i-1], sub[i], sub[i+1], half, style, limit))
	}
	out = append(out, segmentOutline(sub[len(sub)-2], sub[len(sub)-1], half))
	// Each cap reaches away from the line, so the first one looks back the way
	// the line came and the last one carries on the way it was going.
	out = append(out, capOutline(sub[0], unit(sub[0].X-sub[1].X, sub[0].Y-sub[1].Y), half, style.Cap))
	out = append(out, capOutline(sub[len(sub)-1], unit(sub[len(sub)-1].X-sub[len(sub)-2].X, sub[len(sub)-1].Y-sub[len(sub)-2].Y), half, style.Cap))
	return out
}

// segmentOutline is the rectangle of half-width on either side of one straight
// piece, ordered so that it always winds the same way round whatever direction
// the line runs: the normal side first, then across, then back.
func segmentOutline(a, b Point, half float64) []Point {
	d := unit(b.X-a.X, b.Y-a.Y)
	nx, ny := -d.Y*half, d.X*half
	return []Point{
		{a.X + nx, a.Y + ny},
		{b.X + nx, b.Y + ny},
		{b.X - nx, b.Y - ny},
		{a.X - nx, a.Y - ny},
	}
}

// joinOutline is the wedge outside the corner at b, where the segment from a
// turns into the segment to c. The two rectangles always overlap on the inside
// of the turn, and it is the outside that they fall short of, so the wedge goes
// from the outside corner of one to the outside corner of the other.
func joinOutline(a, b, c Point, half float64, style StrokeStyle, limit float64) []Point {
	in := unit(b.X-a.X, b.Y-a.Y)
	next := unit(c.X-b.X, c.Y-b.Y)
	sweep := in.X*next.Y - in.Y*next.X
	if sweep == 0 {
		// Straight on, or straight back: the rectangles meet edge to edge, or
		// lie on top of each other, and there is no wedge to fill.
		return nil
	}
	// The outside of a turn is a quarter turn from the way the line runs, and
	// which of the two quarter turns it is depends on which way the corner
	// turns: a corner opening to the left leaves a gap on its right, and the
	// other way round. The cross product says which way the corner turned.
	side := 1.0
	if sweep > 0 {
		side = -1
	}
	outer := func(d Point) Point {
		return Point{-d.Y * side, d.X * side}
	}
	p1 := offset(b, outer(in), half)
	p2 := offset(b, outer(next), half)

	// Every wedge is closed back through the corner itself. The two rectangles
	// already meet over the inside of the turn, but the corner is the one point
	// they both stop short of, so without it a join is a notch cut into the
	// band where the two sides should have met.
	switch style.Join {
	case JoinRound:
		// The outer normal of a segment turns by exactly as much as the segment
		// does, so the arc from one to the other spans the corner itself.
		from := outer(in)
		turn := math.Atan2(sweep, in.X*next.X+in.Y*next.Y)
		arc := arcPoints(b, half, math.Atan2(from.Y, from.X), math.Atan2(from.Y, from.X)+turn, 1)
		return append(arc, b)
	case JoinMiter:
		if m, ok := miterPoint(b, p1, p2, in, next, sweep, half, limit); ok {
			return []Point{p1, m, p2, b}
		}
	}
	return []Point{p1, b, p2}
}

// offset is a point a distance along a direction, which is what every corner of
// a stroke is: the corner, plus the normal, times half the width.
func offset(at, dir Point, half float64) Point {
	return Point{at.X + dir.X*half, at.Y + dir.Y*half}
}

// capOutline is what a stroke leaves at the end of a line, sitting on the end
// point and reaching outwards along d, which points away from the line: a butt
// cap adds nothing, a square one extends the rectangle by a half-width to the
// side as well, and a round one bulges a half-disc past the end.
func capOutline(at, d Point, half float64, cap LineCap) []Point {
	switch cap {
	case CapRound:
		// The half-disc is the side of the end that the line does not already
		// cover: from one side of the band, round through d, to the other. That
		// is the way round the far side, so the sweep is the short one.
		n := Point{-d.Y * half, d.X * half}
		return arcPoints(at, half, math.Atan2(n.Y, n.X), math.Atan2(-n.Y, -n.X), -1)
	case CapSquare:
		// The two outer corners past the end, which between the band on the
		// far side of the end point is all the cap there is to add.
		n := Point{-d.Y * half, d.X * half}
		return []Point{
			{at.X + n.X, at.Y + n.Y},
			{at.X - n.X, at.Y - n.Y},
			{at.X - n.X + d.X*half, at.Y - n.Y + d.Y*half},
			{at.X + n.X + d.X*half, at.Y + n.Y + d.Y*half},
		}
	}
	return nil
}

// capDisc is the round cap of a stroke that has no line to cap.
func capDisc(at Point, half float64) []Point {
	return arcPoints(at, half, 0, 2*math.Pi, 1)
}

// miterPoint is where the outer edge of the incoming segment meets the outer
// edge of the outgoing one, and whether they meet close enough to keep. Both
// edges run along their own segment, half a width out from it, so the corner
// between them is found by crossing those two lines — and the spike is only
// worth drawing while it stays within the miter limit, which counts the tip's
// distance from the corner in half-widths, the way SVG measures it.
func miterPoint(b, p1, p2, in, next Point, sweep, half, limit float64) (Point, bool) {
	if sweep == 0 {
		return Point{}, false
	}
	// p1 + t*in = p2 + s*next, solved for t by crossing both sides with next.
	t := ((p2.X-p1.X)*next.Y - (p2.Y-p1.Y)*next.X) / sweep
	m := Point{p1.X + in.X*t, p1.Y + in.Y*t}
	if math.Hypot(m.X-b.X, m.Y-b.Y) > limit*half {
		return Point{}, false
	}
	return m, true
}

// arcPoints walks the circle of radius r around c from one angle to the other,
// turning the way step says, which for a join is the way the corner opens and
// for a cap is the way out from the line. A rounding is nothing but this walk,
// which is why it is a list of points and not a shape of its own.
func arcPoints(c Point, r, from, to, step float64) []Point {
	turn := to - from
	if step < 0 && turn > 0 {
		turn -= 2 * math.Pi
	}
	if step > 0 && turn < 0 {
		turn += 2 * math.Pi
	}
	// The gap between a chord and the arc it stands in for shrinks with the
	// square of the number of chords, so a long way round a big circle wants
	// proportionally more of them than a short way round a small one.
	n := max(int(math.Ceil(1.2*r*math.Abs(turn))), 2)
	pts := make([]Point, 0, n+1)
	for i := range n + 1 {
		angle := from + turn*float64(i)/float64(n)
		pts = append(pts, Point{c.X + r*math.Cos(angle), c.Y + r*math.Sin(angle)})
	}
	return pts
}

// wound returns a piece ordered so that it comes out the same way round as the
// rectangles it overlaps, since a piece that came out the other way would take
// its overlap away instead of adding to it.
func wound(pts []Point) []Point {
	if len(pts) > 2 && turn(pts) > 0 {
		slices.Reverse(pts)
	}
	return pts
}

// turn is twice the area a piece encloses, positive when it is wound one way
// round and negative the other. A point at a time in a list of at most a few
// dozen is far cheaper than tracking the winding of the whole outline.
func turn(pts []Point) float64 {
	sum := 0.0
	for i, p := range pts {
		q := pts[(i+1)%len(pts)]
		sum += p.X*q.Y - q.X*p.Y
	}
	return sum
}

// unit is a vector of the given length scaled to one, and the zero vector when
// it has no length at all.
func unit(x, y float64) Point {
	l := math.Hypot(x, y)
	if l == 0 {
		return Point{}
	}
	return Point{x / l, y / l}
}
