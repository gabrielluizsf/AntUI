package css

import "strings"

// gridBackwardName splits a leading "-" off a line name, the way CSS counts a named line from the end of the axis.
// gridBackwardName splits a leading "-" off a line name, the way CSS counts a
// named line from the end of the axis.
func gridBackwardName(raw string) (string, bool) {
	if strings.HasPrefix(raw, "-") {
		return raw[1:], true
	}
	return raw, false
}
