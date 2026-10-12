package sdk

// Major is the first part, or 0 when there is none.
func (v Version) Major() int {
	if len(v.Parts) == 0 {
		return 0
	}
	return v.Parts[0]
}
