package css

import (
	"strings"
)

// ClassList splits a space-separated set of class names.
func ClassList(s string) []string {
	return strings.Fields(s)
}
