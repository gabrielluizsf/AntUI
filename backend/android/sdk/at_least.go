package sdk

// AtLeast reports whether v is w or newer.
func (v Version) AtLeast(w Version) bool { return v.Compare(w) >= 0 }
