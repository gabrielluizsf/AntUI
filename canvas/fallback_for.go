package canvas

// fallbackFor is the face of the chain that has a shape for r, past this one,
// and nil when none of them does. Eight is further down than any chain is
// built, so one written round in a circle stops rather than running forever.
func (f *Face) fallbackFor(r rune) *Face {
	for i, n := 0, f.next; n != nil && i < 8; i, n = i+1, n.next {
		if n.file != nil && n.file.Cmap.Glyph(r) != 0 {
			return n
		}
	}
	return nil
}
