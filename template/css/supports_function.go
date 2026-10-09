package css

import (
	"strings"
)

// function answers one of the condition's own functions. Anything else is a
// function the reader has no rule for, which leaves the whole condition
// unreadable rather than guessed at.
func (r *supportsReader) function(name, inner string) bool {
	switch strings.ToLower(name) {
	case "selector":
		return supportsSelector(inner)
	case "font-format":
		// Only TrueType outlines are read from a file, under either name
		// the spec gives them; the woff containers and the CFF outlines of
		// an OpenType file are formats no @font-face here can open.
		f := unquote(strings.ToLower(strings.TrimSpace(inner)))
		return f == "truetype" || f == "ttf"
	case "font-tech":
		// One technology is promised: the colour bitmaps of a CBDT
		// table, which an @font-face reads and every frame draws. The
		// rest — the COLR, SVG and sbix colour formats, variations,
		// layout features — are not read, so no stylesheet is told they
		// are.
		f := unquote(strings.ToLower(strings.TrimSpace(inner)))
		return f == "color-cbdt"
	}
	r.bad = true
	return false
}
