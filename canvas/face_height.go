package canvas

// Height is how far apart two lines of this face are, in pixels.
func (f *Face) Height() int {
	if f == nil {
		return builtin.height
	}
	return f.height
}
