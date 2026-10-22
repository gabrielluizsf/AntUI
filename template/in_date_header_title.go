package template

// inDateHeaderTitle reports whether (mx,my) is in the header row between the
// two arrow gutters — the month and the year, not the arrows beside them —
// which is what a tap toggles the month grid on.
func inDateHeaderTitle(x, y, u, mx, my int) bool {
	header := dateHeaderHeight(u)
	if my < y || my >= y+header {
		return false
	}
	gut := dateArrowGutter(u)
	w, _ := datePickerSize(u)
	return mx >= x+gut && mx < x+w-gut
}
