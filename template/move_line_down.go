package template

import "github.com/gabrielluizsf/antui/canvas"

func moveLineDown(text string, lines [][2]int, li, in int) int {
	if li >= len(lines)-1 {
		return lines[len(lines)-1][1]
	}
	want := canvas.TextWidth(text[lines[li][0] : lines[li][0]+in])
	ls, le := lines[li+1][0], lines[li+1][1]
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
