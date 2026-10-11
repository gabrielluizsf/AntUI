package svg

import "math"

// copyRange is which copies of a picture written over [lo, hi) — in the
// coordinates of the tile, however far out those run — reach into the tile
// itself, [0, size): the copy at the offset n*size, for each n from the first
// to the last of them. ok answers false where no copy reaches in at all, and
// also where they reach so far out that there would be too many of them to
// draw one by one — at most five along one way, which is a picture spanning
// about five tiles, and anything beyond that takes the plain cut instead.
func copyRange(lo, hi, size float64) (first, last int, ok bool) {
	f := math.Floor(-hi/size) + 1
	l := math.Ceil((size-lo)/size) - 1
	if f > l || l-f >= 5 || math.Abs(f) > 1e6 || math.Abs(l) > 1e6 {
		return 0, 0, false
	}
	return int(f), int(l), true
}
