package css

import (
	"strings"
)

// clauseWord reads the layer or supports keyword an @import may carry after
// its URL, or "" when what comes next is something else — a media type, or
// the end of the statement. s is already trimmed.
func clauseWord(s string) string {
	for _, clause := range []string{"layer", "supports"} {
		if len(s) < len(clause) || !strings.EqualFold(s[:len(clause)], clause) {
			continue
		}
		if len(s) == len(clause) {
			return clause
		}
		switch s[len(clause)] {
		case '(', ' ', '\t', '\r', '\n':
			return clause
		}
	}
	return ""
}
