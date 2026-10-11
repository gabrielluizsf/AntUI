package backend

// Has reports whether every modifier in want is held.
func (m Mod) Has(want Mod) bool { return m&want == want }
