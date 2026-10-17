package svg

import "strings"

// hasClass says whether an element's class attribute names a class, read as
// the whitespace-separated list CSS reads it as.
func hasClass(class, name string) bool {
	for _, f := range strings.Fields(class) {
		if f == name {
			return true
		}
	}
	return false
}
