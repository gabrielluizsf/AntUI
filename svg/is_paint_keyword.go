package svg

import "strings"

// isPaintKeyword says whether a paint is the given CSS keyword. The comparison
// is case-insensitive and ignores surrounding space, as reading a keyword does
// everywhere else here.
func isPaintKeyword(raw, keyword string) bool {
	return strings.EqualFold(strings.TrimSpace(raw), keyword)
}
