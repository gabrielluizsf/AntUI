package svg

import "strings"

// declAttrs takes the name="value" pairs out of an xml declaration. The
// declaration is not a tag and its words are only these few, so a small
// scanner of its own reads them rather than bending the tag reader to it.
func declAttrs(s string) map[string]string {
	said := map[string]string{}
	for {
		s = strings.TrimLeft(s, " \t\r\n")
		if s == "" {
			return said
		}
		i := strings.IndexAny(s, "= \t\r\n")
		if i <= 0 {
			return said
		}
		name := strings.ToLower(s[:i])
		s = strings.TrimLeft(s[i:], " \t\r\n=")
		if s == "" {
			return said
		}
		quote := s[0]
		if quote != '"' && quote != '\'' {
			return said
		}
		s = s[1:]
		j := strings.IndexByte(s, quote)
		if j < 0 {
			return said
		}
		said[name] = s[:j]
		s = s[j+1:]
	}
}
