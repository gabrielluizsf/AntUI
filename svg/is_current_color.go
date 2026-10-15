package svg

// isCurrentColor says whether a paint is the keyword that stands for the colour
// the drawing is being painted in. A drawing that spells it `currentcolor` means
// the same thing as one that does not.
func isCurrentColor(s string) bool {
	return isPaintKeyword(s, "currentcolor")
}
