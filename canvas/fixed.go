package canvas

// Fixed reports whether every character is the same width.
func (f *Face) Fixed() bool { return f == nil || f.file == nil }
