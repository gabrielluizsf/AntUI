package css

import (
	"strconv"
	"strings"
)

// parseColumnCount reads a column count. A whole number is the count; auto and
// an unreadable number leave it to the width, as the initial value does.
func parseColumnCount(raw string) (int, bool) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "auto" {
		return 0, true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}
