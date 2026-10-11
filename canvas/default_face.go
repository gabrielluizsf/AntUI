package canvas

// DefaultFace is what is being used now, and is never nil.
func DefaultFace() *Face {
	if f := defaultFace.Load(); f != nil {
		return f
	}
	return builtin
}
