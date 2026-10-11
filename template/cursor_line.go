package template

// cursorLine finds which visual line a byte cursor sits on and its offset
// within that line.
func cursorLine(lines [][2]int, cursor int) (li, in int) {
	for i, l := range lines {
		if cursor >= l[0] && cursor <= l[1] {
			return i, cursor - l[0]
		}
	}
	if len(lines) == 0 {
		return 0, 0
	}
	last := lines[len(lines)-1]
	return len(lines) - 1, last[1] - last[0]
}
