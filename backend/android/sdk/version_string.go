package sdk

import (
	"strconv"
	"strings"
)

// String writes the version back out the way it came in.
func (v Version) String() string {
	if len(v.Parts) == 0 {
		return ""
	}
	var b strings.Builder
	for i, n := range v.Parts {
		if i > 0 {
			b.WriteByte('.')
		}
		b.WriteString(strconv.Itoa(n))
	}
	if v.Pre != "" {
		b.WriteByte('-')
		b.WriteString(v.Pre)
	}
	return b.String()
}
