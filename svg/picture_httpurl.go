package svg

import "strings"

// pictureHTTPURL is the address a href is fetched from where it is one at all:
// http and https as they stand, and a protocol-relative `//host/path` taken to
// https, which is the scheme a drawing with no page above it means. Anything
// else — a data URI, a file path — is not an address this function names.
func pictureHTTPURL(href string) (string, bool) {
	switch {
	case strings.HasPrefix(href, "http://"), strings.HasPrefix(href, "https://"):
		return href, true
	case strings.HasPrefix(href, "//"):
		return "https:" + href, true
	}
	return "", false
}
