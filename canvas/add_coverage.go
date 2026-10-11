package canvas

import "math"

// addCoverage adds the length of a stretch of line that falls inside, spread
// over the pixels it falls on. The stretch is clipped to the band, and a pixel
// it starts or stops inside is worth only the part of it the line really
// covers, which is the whole of the smoothing along a near-vertical edge.
func addCoverage(row []float64, x0 int, from, to float64) {
	from, to = max(from, float64(x0)), min(to, float64(x0+len(row)))
	if to <= from {
		return
	}
	first := max(int(math.Floor(from)), x0)
	last := min(int(math.Ceil(to)), x0+len(row))
	for x := first; x < last; x++ {
		row[x-x0] += min(to, float64(x+1)) - max(from, float64(x))
	}
}
