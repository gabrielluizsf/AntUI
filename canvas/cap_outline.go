package canvas

import "math"

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
