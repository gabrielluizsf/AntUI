package template

import (
	"unicode/utf8"

	"github.com/gabrielluizsf/antui/canvas"
)

// cutWord clips a word from s to e down to the longest rune prefix whose width
// fits avail, leaving the rest on the next line.
func cutWord(text string, s, e, avail int) int {
	w := 0
	i := s
	for i < e {
		r, sz := utf8.DecodeRuneInString(text[i:])
		gw := canvas.TextWidth(string(r))
		if i > s && w+gw > avail {
			break
		}
		w += gw
		i += sz
	}
	return i
}
