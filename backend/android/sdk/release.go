package sdk

// Release reports whether this is a release rather than a preview. Previews
// are never chosen on their own; see [pick].
func (v Version) Release() bool { return v.Pre == "" }
