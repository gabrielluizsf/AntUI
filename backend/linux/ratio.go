//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// ratio turns a width-to-height ratio into the fraction X11 wants. Two
// decimal places of denominator is enough for every ratio a game uses — 16:9,
// 16:10, 4:3, 21:9 — and asking for more would only invite a rounding that
// makes the window drift by a pixel each drag.
func ratio(aspect float64) (num, den uint32) {
	const scale = 1000
	n := uint32(aspect*scale + 0.5)
	if n == 0 {
		return 1, 1
	}
	return n, scale
}
