package antui

// slop is a threshold in points converted to pixels. A window whose backend
// does not know the display scale gets the number unchanged, which is right
// for a desktop and merely small on a phone that will not say.
func (win *Window) slop(points float64) float64 {
	if win.scale > 0 {
		return points * win.scale
	}
	return points
}
