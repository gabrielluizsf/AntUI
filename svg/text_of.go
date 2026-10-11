package svg

import "strings"

// textOf is all the writing inside an element, put together, with the
// whitespace at the two ends of it taken off, because whitespace around a
// drawing is there to make it readable and says nothing about where it is.
func textOf(e *element) string {
	var out strings.Builder
	for _, p := range e.Parts {
		out.WriteString(p.Text)
	}
	return strings.TrimSpace(out.String())
}
