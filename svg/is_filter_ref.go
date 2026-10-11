package svg

import "strings"

// isFilterRef says the token is a `url(...)` — the reference to a `<filter>`
// element whose id is kept here for the whole drawing to follow.
func isFilterRef(tok string) bool {
	t := strings.TrimSpace(tok)
	return strings.HasPrefix(strings.ToLower(t), "url(") ||
		strings.EqualFold(t, "url")
}
