package svg

import "strings"

// hrefScheme is the scheme an address names, lower-cased and without the colon,
// or "" where it names none — a relative address or a plain path. A scheme is
// the address up to the first colon, starting with a letter and carrying only
// letters, digits, `+`, `-` and `.` the way a URL spells one, so a path that
// happens to hold a colon is not mistaken for an address somewhere else.
func hrefScheme(target string) string {
	i := strings.IndexByte(target, ':')
	if i <= 0 {
		return ""
	}
	scheme := target[:i]
	for j := 0; j < len(scheme); j++ {
		c := scheme[j]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z':
		case j > 0 && (c >= '0' && c <= '9' || c == '+' || c == '-' || c == '.'):
		default:
			return ""
		}
	}
	return strings.ToLower(scheme)
}
