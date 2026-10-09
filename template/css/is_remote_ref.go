package css

import (
	"strings"
)

// isRemoteRef reports whether an @import names something outside this
// machine: another host, a data: document, anything with a scheme behind it.
// Reading a stylesheet means reading a file beside this one; a canvas that
// went out to fetch one would be a network client as well as a renderer.
func isRemoteRef(ref string) bool {
	if strings.HasPrefix(ref, "//") {
		return true
	}
	// A colon before any '/' is a scheme, unless it is a single letter —
	// that is a drive root, and reads as a path like any other.
	if i := strings.Index(ref, ":"); i > 0 && i != 1 && !strings.ContainsAny(ref[:i], `/\`) {
		return true
	}
	return false
}
