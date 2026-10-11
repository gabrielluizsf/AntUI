package template

// textViewTop is the first line the box should show for cursor, so the caret
// stays visible as it moves. lines holds the wrapped layout, visible how many
// rows fit in the box.
func textViewTop(lines [][2]int, cursor, visible int) int {
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	return max(min(viewTop, max(len(lines)-visible, 0)), 0)
}
