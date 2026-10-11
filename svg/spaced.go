package svg

// spaced is how many letters of a piece the drawing asked for room after: a
// line break is the end of a line rather than a letter at the end of one, and
// is left out of the count.
func spaced(st Style, text string) int {
	if st.LetterSpacing == 0 {
		return 0
	}
	n := 0
	for _, r := range text {
		if r != '\n' && r != '\r' {
			n++
		}
	}
	return n
}
