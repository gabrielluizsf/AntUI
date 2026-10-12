package svg

import "strings"

// keepText appends the text gathered in b to e as one part — collapsed unless
// preserve asked for it exactly as written — and clears the builder.
func keepText(e *element, b *strings.Builder, preserve bool) {
	s := b.String()
	if !preserve {
		s = collapseSpace(s)
	}
	if s != "" {
		e.Parts = append(e.Parts, part{Text: s})
	}
	b.Reset()
}
