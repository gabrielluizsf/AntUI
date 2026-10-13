package template

// dateArrowRect is the square the prev (left) or next (right) arrow draws in
// the calendar header, leaving the day grid aligned as one table below it.
func dateArrowRect(x, y, u int, prev bool) (ax, ay, aw int) {
	aw = 8 * u
	header := dateHeaderHeight(u)
	ax, ay = x+4*u, y+(header-aw)/2
	if !prev {
		w, _ := datePickerSize(u)
		ax = x + w - 4*u - aw
	}
	return
}
