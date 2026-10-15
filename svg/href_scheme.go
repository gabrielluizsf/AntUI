package svg

import (
	"strings"
)

// hrefScheme is the scheme of an address, lower-cased: what comes before the
// first colon, and only where it is spelled the way a URL spells one.
func hrefScheme(target string) string {
	i := strings.IndexByte(target, ':')
	if i <= 0 {
		return ""
	}

	// A Windows drive letter followed by a slash or backslash
	// identifies a local path, not a URL scheme.
	if i == 1 && isASCIIAlpha(target[0]) &&
		len(target) > 2 &&
		(target[2] == '\\' || target[2] == '/') {
		return ""
	}

	scheme := target[:i]
	for j := 0; j < len(scheme); j++ {
		c := scheme[j]
		switch {
		case isASCIIAlpha(c):
		case j > 0 && (c >= '0' && c <= '9' ||
			c == '+' || c == '-' || c == '.'):
		default:
			return ""
		}
	}

	return strings.ToLower(scheme)
}

func isASCIIAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z'
}
