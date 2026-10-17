package svg

import "fmt"

// readAttrValue reads the value an attribute selector compares: a quoted
// string, or a bare name running up to the closing bracket or a space.
func readAttrValue(s string, i int) (string, int, error) {
	if i < len(s) && (s[i] == '"' || s[i] == '\'') {
		q := s[i]
		j := i + 1
		for j < len(s) && s[j] != q {
			j++
		}
		if j >= len(s) {
			return "", i, fmt.Errorf("a quoted attribute value is never closed")
		}
		return s[i+1 : j], j + 1, nil
	}
	j := i
	for j < len(s) && s[j] != ']' && !isSpace(s[j]) {
		j++
	}
	if j == i {
		return "", i, fmt.Errorf("an attribute selector has no value")
	}
	return s[i:j], j, nil
}
