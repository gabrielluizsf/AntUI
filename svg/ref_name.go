package svg

import (
	"fmt"
	"strings"
)

// refName reads the name between the % or & at s[i] and the `;` that closes
// it: the name itself, and how many bytes the whole reference took. A
// reference without its `;`, or a name XML does not allow, is an error rather
// than letters a reader would guess a meaning for.
func refName(s string, i int) (string, int, error) {
	rest := s[i+1:]
	end := strings.IndexByte(rest, ';')
	if end < 0 {
		return "", 0, fmt.Errorf("the name after a `%c` is never closed with a `;`", s[i])
	}
	name := rest[:end]
	if !validName(name) {
		return "", 0, fmt.Errorf("`%c%s;` is not a name XML allows", s[i], name)
	}
	return name, end + 2, nil
}
