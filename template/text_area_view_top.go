package template

// textAreaViewTop scrolls the field up so the line the cursor is on stays
// visible at the bottom edge.
func textAreaViewTop(cursor int, lines [][2]int, visible int) int {
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	return max(min(viewTop, len(lines)-visible), 0)
}
