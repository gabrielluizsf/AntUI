package template

import "unicode/utf8"

// scanWord advances past one word, returning where it started and where it
// ends. A word runs to the next soft space or hard newline.
func scanWord(text string, i int) (start, end int) {
	start = i
	for i < len(text) && text[i] != '\n' && !isSoftSpace(text[i]) {
		_, sz := utf8.DecodeRuneInString(text[i:])
		i += sz
	}
	return start, i
}
