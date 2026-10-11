package svg

import (
	"strconv"
	"strings"
)

// numbers reads a list of numbers the way a filter writes one: separated by
// spaces, commas or both, and all of them the same kind of thing. A piece that
// is not a number answers false for the whole list — the caller wanted a list
// and what it has is not one.
func numbers(s string) ([]float64, bool) {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == ','
	})
	out := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}
