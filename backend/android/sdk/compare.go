package sdk

import "cmp"

// Compare orders two versions: by their numbers, then by whether they are a
// release, so 37.0.0 sorts above 37.0.0-rc2. A version with fewer parts is
// the smaller one when the parts they share are equal, so 37.0 is below
// 37.0.1.
func (v Version) Compare(w Version) int {
	for i := range max(len(v.Parts), len(w.Parts)) {
		var a, b int
		if i < len(v.Parts) {
			a = v.Parts[i]
		}
		if i < len(w.Parts) {
			b = w.Parts[i]
		}
		if c := cmp.Compare(a, b); c != 0 {
			return c
		}
	}
	switch {
	case v.Pre == w.Pre:
		return 0
	case v.Pre == "":
		return 1
	case w.Pre == "":
		return -1
	}
	return cmp.Compare(v.Pre, w.Pre)
}
