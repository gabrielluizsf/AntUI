package svg

import "strings"

// isCurrentColor says whether a paint is the keyword that stands for the colour
// the drawing is being painted in. The comparison is case-insensitive because
// CSS keywords are, and a drawing that spells it `currentcolor` means the same
// thing as one that does not.
func isCurrentColor(s string) bool {
	return strings.EqualFold(strings.TrimSpace(s), "currentcolor")
}
