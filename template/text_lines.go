package template

import (
	"unicode/utf8"

	"github.com/gabrielluizsf/antui/canvas"
)

// textLines breaks multiline text into wrapped visual lines. Each entry is
// the [start, end) byte range of one visual line; a '\n' ends a line and is
// not part of it. maxRef is the wrap width at reference scale.
func textLines(text string, maxRef int) [][2]int {
	var lines [][2]int
	i := 0
	for i <= len(text) {
		start := i
		w := 0
		brk := -1
		for i < len(text) {
			if text[i] == '\n' {
				brk = i
				break
			}
			r, size := utf8.DecodeRuneInString(text[i:])
			gw := canvas.TextWidth(string(r))
			if w+gw > maxRef && w > 0 {
				brk = i
				break
			}
			w += gw
			i += size
		}
		if brk == -1 {
			lines = append(lines, [2]int{start, len(text)})
			return lines
		}
		lines = append(lines, [2]int{start, brk})
		if text[brk] == '\n' {
			i = brk + 1
		} else {
			i = brk
		}
	}
	return lines
}
