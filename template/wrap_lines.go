package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// wrapLines breaks text into visual lines like textLines, but honours the
// white-space and overflow-wrap semantics the engine models:
//
//   - normal and pre-line drop leading whitespace from each line and count a
//     run of spaces as one width;
//   - pre and pre-wrap keep every space as written and pre-wrap alone soft
//     wraps;
//   - nowrap and pre never soft wrap, so an overlong line spills past maxRef;
//   - overflow-wrap break-word/anywhere cuts a word that overflows the line;
//     normal lets the word spill onto its own overlong line.
func wrapLines(text string, maxRef int, ws, ow uint8) [][2]int {
	softWrap := ws != css.WhiteSpaceNowrap && ws != css.WhiteSpacePre
	breakWord := ow == css.OverflowWrapBreakWord || ow == css.OverflowWrapAnywhere
	keepSpaces := ws == css.WhiteSpacePre || ws == css.WhiteSpacePreWrap
	spaceW := canvas.TextWidth(" ")

	var lines [][2]int
	for i := 0; i < len(text); {
		if !keepSpaces {
			i = skipSoftSpaces(text, i)
			if i >= len(text) {
				break
			}
		}
		end, next := wrapLine(text, i, maxRef, spaceW, keepSpaces, softWrap, breakWord)
		if end == len(text) && end == i {
			break
		}
		lines = append(lines, [2]int{i, end})
		i = next
	}
	return lines
}
