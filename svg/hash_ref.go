package svg

import "strings"

// hashRef reads the `#id` out of either spelling of a reference: the plain one
// an href is written in, and the `url(#id)` a paint takes.
func hashRef(s string) (string, bool) {
	if id, ok := paintRef(s); ok {
		return id, true
	}
	s = strings.Trim(strings.TrimSpace(s), "\"'")
	if !strings.HasPrefix(s, "#") {
		return "", false
	}
	return s[1:], s[1:] != ""
}
