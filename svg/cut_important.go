package svg

import "strings"

// cutImportant takes a trailing `!important` off a declaration's value and
// says whether there was one, allowing for the space CSS allows between the
// `!` and the word.
func cutImportant(v string) (value string, important bool) {
	s := strings.TrimRight(v, " \t\r\n")
	i := strings.LastIndexByte(s, '!')
	if i < 0 || strings.TrimSpace(strings.ToLower(s[i+1:])) != "important" {
		return v, false
	}
	return strings.TrimSpace(s[:i]), true
}
