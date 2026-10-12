package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// wrapLine scans one visual line from start and returns where the line ends
// and where the scan resumes, which is past a hard newline when one ended it.
func wrapLine(text string, start, maxRef, spaceW int, keepSpaces, softWrap, breakWord bool) (end, next int) {
	i := start
	w := 0
	brk := -1
	for i < len(text) && text[i] != '\n' {
		if isSoftSpace(text[i]) {
			var runW int
			runW, i = spaceRun(text, i, keepSpaces, spaceW)
			// A run after real content is always a break opportunity:
			// ending the line there drops the trailing spaces naturally,
			// even when the run itself would overflow a hair.
			if w > 0 {
				brk = i
			}
			if w+runW <= maxRef {
				w += runW
			}
			continue
		}
		wordStart, wordEnd := scanWord(text, i)
		wordW := canvas.TextWidth(text[wordStart:wordEnd])
		i = wordEnd
		if softWrap && w+wordW > maxRef {
			if cut, ok := wrapWord(text, wordStart, wordEnd, wordW, brk, w, maxRef, breakWord); ok {
				return cut, cut
			}
		}
		w += wordW
	}
	end = i
	next = i
	if i < len(text) && text[i] == '\n' {
		next++
	}
	return end, next
}
