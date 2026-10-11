package svg

import (
	"strings"
	"unicode/utf8"
)

// collapseSpace turns every run of whitespace into a single space, so that text
// wrapped across lines in the file is the same text as text written on one.
// The spaces at the two ends are kept: a space at the edge of a piece of
// writing is what sits between it and the piece written next to it, and only
// whoever reads all of it together knows whether the writing wants it, which
// is what taking them off the whole of it in [textOf] does.
func collapseSpace(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	space := false
	for i := 0; i < len(s); {
		c := s[i]
		if isSpace(c) {
			space = true
			i++
			continue
		}
		if space {
			out.WriteByte(' ')
		}
		space = false
		_, size := utf8.DecodeRuneInString(s[i:])
		out.WriteString(s[i : i+size])
		i += size
	}
	if space {
		out.WriteByte(' ')
	}
	return out.String()
}
