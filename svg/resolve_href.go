package svg

import (
	"net/url"
	"strings"
)

// resolveHref is the address an href points at once the base the drawing sits
// at is taken into account: a relative address is completed against the base,
// the way XML Base says, and an address that names a scheme of its own is
// already whole. With no base — a drawing read from a string rather than a file
// — the address is left exactly as it was written, so a file beside the drawing
// keeps working while it is being made. An address that will not parse is left
// alone too, to be refused where it is used.
//
// A protocol-relative `//host/path` names no scheme: it takes the base's where
// there is one to take it from, and is otherwise the https a drawing with no
// page above it means.
//
// A base that names no scheme (or that starts with a Windows drive letter) is a
// path on this machine, not a URL, and is joined with the path package rather
// than parsed as a URL: `C:\dir\` parsed as a URL would read `C` as a scheme.
func resolveHref(base, href string) string {
	if href == "" || isWindowsAbsPath(href) {
		return href
	}
	if strings.HasPrefix(href, "//") {
		if scheme := hrefScheme(base); scheme != "" {
			return scheme + ":" + href
		}
		return "https:" + href
	}
	// No base to complete it against, or an address that is already whole.
	if base == "" || hrefScheme(href) != "" {
		return href
	}

	if hrefScheme(base) == "" {
		return joinLocal(base, href)
	}

	b, err := url.Parse(base)
	if err != nil {
		return href
	}
	r, err := url.Parse(href)
	if err != nil {
		return href
	}
	return b.ResolveReference(r).String()
}

func isWindowsAbsPath(path string) bool {
	return len(path) >= 3 &&
		isASCIIAlpha(path[0]) &&
		path[1] == ':' &&
		(path[2] == '\\' || path[2] == '/')
}
