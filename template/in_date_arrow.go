package template

// inDateArrow reports whether (mx,my) is on the left (prev, true) or right
// (next, false) arrow of the header. The arrow is a square under the header
// text baseline, at either end of the row.
func inDateArrow(x, y, u, mx, my int, prev bool) bool {
	ax, aw := x+2*u, 8*u
	if !prev {
		w, _ := datePickerSize(u)
		ax = x + w - 2*u - aw
	}
	return mx >= ax && mx < ax+aw && my >= y && my < y+dateHeaderHeight(u)
}
