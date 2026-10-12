package sdk

import (
	"strconv"
	"strings"
)

// ParseVersion reads a version. It never fails: anything it cannot make sense
// of contributes nothing, so a directory named by something other than a
// version sorts below every real one instead of stopping the search.
func ParseVersion(s string) Version {
	var v Version
	if i := strings.IndexByte(s, '-'); i >= 0 {
		v.Pre, s = s[i+1:], s[:i]
	}
	for _, part := range strings.Split(s, ".") {
		n, err := strconv.Atoi(part)
		if err != nil {
			break
		}
		v.Parts = append(v.Parts, n)
	}
	return v
}
