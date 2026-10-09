package css

import "sort"

// sortCandidates orders the winning rules by specificity and then source
// order, the way the cascade weighs two declarations against each other.
func sortCandidates(got []candidate) {
	sort.SliceStable(got, func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if got[i].spec[k] != got[j].spec[k] {
				return got[i].spec[k] < got[j].spec[k]
			}
		}
		return false
	})
}
