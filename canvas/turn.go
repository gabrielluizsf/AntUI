package canvas

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
