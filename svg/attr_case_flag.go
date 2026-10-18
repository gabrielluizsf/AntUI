package svg

import "fmt"

// attrCaseFlag reads the `i` or `s` that may follow an attribute selector's
// value — `i` compares without case, `s` keeps the case, which is the default
// — and answers where the closing bracket then stands.
func attrCaseFlag(s string, i int) (bool, int, error) {
	if i >= len(s) || (s[i] != 'i' && s[i] != 'I' && s[i] != 's' && s[i] != 'S') {
		return false, i, nil
	}
	j := skipSelSpace(s, i+1)
	if j >= len(s) || s[j] != ']' {
		return false, i, fmt.Errorf("an attribute selector flag is not one CSS names")
	}
	return s[i] == 'i' || s[i] == 'I', j, nil
}
