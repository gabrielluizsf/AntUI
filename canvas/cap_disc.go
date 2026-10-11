package canvas

import "math"

// capDisc is the round cap of a stroke that has no line to cap.
func capDisc(at Point, half float64) []Point {
	return arcPoints(at, half, 0, 2*math.Pi, 1)
}
