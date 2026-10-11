package canvas

// Ascent is how far the tallest letters reach above the baseline, and
// Descent how far the tails go below it.
func (f *Face) Ascent() int {
	if f == nil {
		return builtin.ascent
	}
	return f.ascent
}
