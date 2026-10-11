package svg

import "strings"

// parseViewBox reads the four numbers of a viewBox, answering false for one
// that is not four numbers or that has nothing to show.
func parseViewBox(s string) ([4]float64, bool) {
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\n' || r == '\r' })
	if len(fields) != 4 {
		return [4]float64{}, false
	}
	var v [4]float64
	for i, f := range fields {
		n, err := parseNumber(f)
		if err != nil {
			return [4]float64{}, false
		}
		v[i] = n
	}
	if v[2] <= 0 || v[3] <= 0 {
		return [4]float64{}, false
	}
	return v, true
}
