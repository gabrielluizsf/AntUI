package svg

import "strings"

// readLetterSpacing reads the room left after every letter, in the drawing's
// own units. `normal` and no letter-spacing at all are the same answer, and so
// is a spacing this cannot read: the drawing keeps the letters where they were
// rather than being pulled together by a number that is not one.
func readLetterSpacing(raw string, warn func(string, ...any)) float64 {
	s := strings.TrimSpace(raw)
	if s == "" || strings.EqualFold(s, "normal") {
		return 0
	}
	if v, ok := parseLength(s); ok && v >= 0 {
		return v
	}
	warn("the letter spacing %q is not a length, so the letters keep the room they had", raw)
	return 0
}
