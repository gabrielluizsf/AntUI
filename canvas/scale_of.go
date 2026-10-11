package canvas

// scaleOf is how many device pixels one cell of the built-in font comes out
// at: one for the font as it stands, and whatever a size asked for makes it.
func (f *Face) scaleOf() float64 {
	if f == nil || f.unit <= 0 {
		return 1
	}
	return f.unit
}
