package canvas

import "math"

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
