package canvas

// SetDefaultFace makes every Text, TextWidth and widget in the program use
// this face. nil goes back to the built-in.
func SetDefaultFace(f *Face) { defaultFace.Store(f) }
