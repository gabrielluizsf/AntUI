package svg

// matchesNth says whether a 1-based position is one `An+B` names, for some k of
// zero or more.
func matchesNth(a, b, i int) bool {
	if a == 0 {
		return i == b
	}
	d := i - b
	if d%a != 0 {
		return false
	}
	return d/a >= 0
}
