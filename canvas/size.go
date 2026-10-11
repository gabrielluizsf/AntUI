package canvas

// Size is the face's size in pixels per em, and 0 for the built-in.
func (f *Face) Size() float64 {
	if f == nil || f.file == nil {
		return 0
	}
	return f.size
}
