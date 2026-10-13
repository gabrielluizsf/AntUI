package backend

import "math"

// clamp reduces a size to what the limits allow. It is the whole of the
// arithmetic in this package and is deliberately free of any window: it is
// what the tests can check on every platform, including the ones with no
// window system at all.
func (l Limits) Clamp(width, height int) (int, int) {
	width, height = max(width, 1), max(height, 1)
	width = bound(width, l.MinWidth, l.MaxWidth)
	height = bound(height, l.MinHeight, l.MaxHeight)

	if l.Aspect > 0 {
		// The height follows the width, because a window is nearly always
		// dragged by a side or a corner and the width is the axis a person
		// thinks in. When the height that follows is out of bounds the width
		// follows the height instead, so a ratio can never push a window past
		// a limit it was also given — the limits are the harder promise.
		wanted := max(int(math.Round(float64(width)/l.Aspect)), 1)
		if held := bound(wanted, l.MinHeight, l.MaxHeight); held == wanted {
			height = wanted
		} else {
			height = held
			width = bound(max(int(math.Round(float64(held)*l.Aspect)), 1),
				l.MinWidth, l.MaxWidth)
		}
	}
	return width, height
}
