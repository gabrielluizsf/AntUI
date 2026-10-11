package svg

import "math"

// modUnit is how far along a tile a point is, wrapped into [0, size): the
// subtraction that makes a pattern repeat, done so that a point before the
// first tile lands in the last one rather than nowhere. A tile is never no
// size here — a tile that was has already said so and paints nothing.
func modUnit(a, size float64) float64 {
	v := math.Mod(a, size)
	if v < 0 {
		v += size
	}
	return v
}
