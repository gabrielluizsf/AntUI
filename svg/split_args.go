package svg

import "strings"

// splitArgs breaks the inside of a transform function into its numbers, which
// may be separated by commas, by spaces, or by both.
func splitArgs(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
}
