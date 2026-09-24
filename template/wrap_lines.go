package template

// Wrapping text into lines honouring white-space and overflow-wrap. The
// engine measures at the reference scale (canvas.TextWidth), the same ruler
// textLines uses, so maxRef is the wrap width in reference units.

import (
	"unicode/utf8"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// isSoftSpace reports whether c is a space the white-space handling manages.
func isSoftSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\f'
}

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
outer:
	for i := 0; i < len(text); {
		if !keepSpaces {
			for i < len(text) && isSoftSpace(text[i]) {
				i++
			}
			if i >= len(text) {
				break
			}
		}
		start := i
		w := 0
		brk := -1
		for i < len(text) && text[i] != '\n' {
			if isSoftSpace(text[i]) {
				j := i
				runW := 0
				for j < len(text) && text[j] != '\n' && isSoftSpace(text[j]) {
					if keepSpaces {
						runW += spaceW
					}
					j++
				}
				if !keepSpaces {
					runW = spaceW
				}
				// A run after real content is always a break opportunity:
				// ending the line there drops the trailing spaces naturally,
				// even when the run itself would overflow a hair.
				if w > 0 {
					brk = j
				}
				if w+runW <= maxRef {
					w += runW
				}
				i = j
				continue
			}
			wordStart := i
			for i < len(text) && text[i] != '\n' && !isSoftSpace(text[i]) {
				_, sz := utf8.DecodeRuneInString(text[i:])
				i += sz
			}
			wordW := canvas.TextWidth(text[wordStart:i])
			if softWrap && w+wordW > maxRef {
				if brk != -1 {
					lines = append(lines, [2]int{start, brk})
					i = brk
					continue outer
				}
				if breakWord {
					cut := cutWord(text, wordStart, i, maxRef-w)
					lines = append(lines, [2]int{start, cut})
					i = cut
					continue outer
				}
				if w == 0 {
					// No prior content and no wrap point: the overlong word
					// fills this line and the scan keeps going.
					w += wordW
					continue
				}
				// Prior content but no fitting break: the word spills onto its
				// own overlong line, as CSS overflow-wrap normal allows.
				lines = append(lines, [2]int{start, i})
				continue outer
			}
			w += wordW
		}
		end := i
		if end == len(text) && end == start {
			break
		}
		lines = append(lines, [2]int{start, end})
		if i < len(text) && text[i] == '\n' {
			i++
		}
	}
	return lines
}

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
