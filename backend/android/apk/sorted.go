package apk

import (
	"sort"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// sorted puts the ABIs of a map in a fixed order, so that two builds of the
// same source produce the same archive.
func sorted(m map[sdk.ABI]string) []sdk.ABI {
	out := make([]sdk.ABI, 0, len(m))
	for abi := range m {
		out = append(out, abi)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
