package template

import "github.com/gabrielluizsf/antui/canvas"

// moveLineUp moves the cursor one visual line up, and moveLineDown one down,
// keeping the horizontal position as close as possible.
func moveLineUp(text string, lines [][2]int, li, in int) int {
	if li == 0 {
		return lines[0][0]
	}
	want := canvas.TextWidth(text[lines[li][0] : lines[li][0]+in])
	ls, le := lines[li-1][0], lines[li-1][1]
	best := ls
	for i := ls; i < le; i = utf8NextCursor(text, i) {
		n := utf8NextCursor(text, i)
		w := canvas.TextWidth(text[ls:n])
		if w <= want {
			best = n
		} else {
			break
		}
	}
	return best
}
