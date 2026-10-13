package template

// inDateArrow reports whether (mx,my) is in the left or right arrow area of
// the header. The whole gutter is clickable, not only the small SVG square.
func inDateArrow(x, y, u, mx, my int, prev bool) bool {
	header := dateHeaderHeight(u)
	if my < y || my >= y+header {
		return false
	}
	gut := dateArrowGutter(u)
	if prev {
		return mx >= x && mx < x+gut
	}
	w, _ := datePickerSize(u)
	return mx >= x+w-gut && mx < x+w
}
