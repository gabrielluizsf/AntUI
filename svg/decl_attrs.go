package svg

import (
	"fmt"
	"strings"
)

// declAttrs takes the name="value" pairs out of an xml declaration. The
// declaration is not a tag and its words are only these few, so a small scanner
// of its own reads them rather than bending the tag reader to it. A word other
// than the three the grammar allows — in another case, said twice, or out of
// order — is an error: a declaration that says something XML gives it no word
// for is not a declaration a reader can take.
func declAttrs(s string) (map[string]string, error) {
	said := map[string]string{}
	at := 0
	for {
		s = strings.TrimLeft(s, " \t\r\n")
		if s == "" {
			return said, nil
		}
		i := strings.IndexAny(s, "= \t\r\n")
		if i <= 0 {
			return nil, fmt.Errorf("antui/svg: the xml declaration holds %q, where a word of it was expected", s)
		}
		name := s[:i]
		next := declWord(name)
		if next < at {
			return nil, fmt.Errorf("antui/svg: the xml declaration says %q, which is not a word it may hold here", name)
		}
		at = next + 1
		s = strings.TrimLeft(s[i:], " \t\r\n=")
		if s == "" {
			return nil, fmt.Errorf("antui/svg: the xml declaration gives %s no value", name)
		}
		quote := s[0]
		if quote != '"' && quote != '\'' {
			return nil, fmt.Errorf("antui/svg: the value of %s in the xml declaration is not in quotes", name)
		}
		s = s[1:]
		j := strings.IndexByte(s, quote)
		if j < 0 {
			return nil, fmt.Errorf("antui/svg: the value of %s in the xml declaration is never closed", name)
		}
		said[name] = s[:j]
		s = s[j+1:]
	}
}
