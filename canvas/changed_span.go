package canvas

import "slices"

// changedSpan is the first and the last pixel of two equal-length runs that
// differ, and whether they differ at all. The whole-run comparison is one
// memcmp; the narrowing walk only runs on a run that moved, and then usually
// stops at the first pixel.
func changedSpan(a, b []Color) (int, int, bool) {
	n := min(len(a), len(b))
	if slices.Equal(a[:n], b[:n]) {
		return 0, 0, false
	}
	left := 0
	for a[left] == b[left] {
		left++
	}
	right := n - 1
	for a[right] == b[right] {
		right--
	}
	return left, right, true
}
