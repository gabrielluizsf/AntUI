package canvas

// offset is a point a distance along a direction, which is what every corner of
// a stroke is: the corner, plus the normal, times half the width.
func offset(at, dir Point, half float64) Point {
	return Point{at.X + dir.X*half, at.Y + dir.Y*half}
}
