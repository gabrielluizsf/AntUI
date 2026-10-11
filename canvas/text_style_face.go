package canvas

// face is the face this style draws in: its own when it named one, and the
// one handed in — the canvas's, or the program's — when it did not.
func (o TextStyle) face(fallback *Face) *Face {
	if o.Face != nil {
		return o.Face
	}
	return fallback
}
