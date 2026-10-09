package css

import (
	"fmt"
	"strings"
)

// mediaPx reads a media width like "720px" or "720" into pixels. Fractions
// round; junk returns ok=false and the callers warn.
func mediaPx(s string) (int, bool) {
	t := strings.TrimSpace(strings.ToLower(strings.TrimSuffix(strings.TrimSpace(s), "px")))
	var f float64
	if _, err := fmt.Sscanf(t, "%g", &f); err != nil {
		return 0, false
	}
	return int(f + 0.5), true
}
