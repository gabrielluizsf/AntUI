package template

// inDateGrid reports whether (mx,my) is inside the day grid.
func inDateGrid(x, y, u, mx, my int) bool {
	gx, gy, gw, gh := dateGridRect(x, y, u)
	return mx >= gx && mx < gx+gw && my >= gy && my < gy+gh
}
