package css

import "strings"

// parseRadiiList expands one to four non-negative pixel radii into the four
// corners, in CSS order: top-left, top-right, bottom-right, bottom-left.
func parseRadiiList(raw string) ([4]int, bool) {
	var out [4]int
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]int, len(parts))
	for i, p := range parts {
		l, err := parseLength(p)
		if err != nil || l.u != unitPx || l.value < 0 {
			return out, false
		}
		vals[i] = int(l.value + 0.5)
	}
	fill := expandFour(vals)
	copy(out[:], fill[:])
	return out, true
}

// parseRadiiXY parses a border-radius value, splitting the elliptical form
// "8px / 16px" into horizontal and vertical radii. With no slash the two axes
// are equal.
func parseRadiiXY(raw string) ([4]int, [4]int, bool) {
	head, tail, hasSlash := splitSlash(raw)
	rx, ok := parseRadiiList(head)
	if !ok {
		return rx, rx, false
	}
	ry := rx
	if hasSlash && strings.TrimSpace(tail) != "" {
		if v, ok := parseRadiiList(tail); ok {
			ry = v
		} else {
			return rx, rx, false
		}
	}
	return rx, ry, true
}
