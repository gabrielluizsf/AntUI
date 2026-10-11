//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"net/url"
	"strings"
)

// pathFromURI is the local path a file:// URI names, and false for anything
// else.
func pathFromURI(uri string) (string, bool) {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return "", false
	}
	if !strings.HasPrefix(uri, "file://") {
		// A bare path is what some sources send, and what a person pasting a
		// path into the window means.
		if strings.HasPrefix(uri, "/") {
			return uri, true
		}
		return "", false
	}
	rest := uri[len("file://"):]
	// file://host/path is a file on another machine unless the host is this
	// one, which is written as an empty host or `localhost`.
	if at := strings.IndexByte(rest, '/'); at > 0 {
		if host := rest[:at]; host != "localhost" {
			return "", false
		}
		rest = rest[at:]
	}
	unescaped, err := url.PathUnescape(rest)
	if err != nil {
		return "", false
	}
	return unescaped, true
}
