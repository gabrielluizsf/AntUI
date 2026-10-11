//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"strings"
)

// parseURIList turns what a drag or a clipboard holds into paths on this
// machine.
//
// The format is one URI a line, with comment lines beginning with a hash, and
// the lines end with a carriage return as well as a newline. Anything that is
// not a local file — a web address, a file on another machine — is left out
// rather than handed on as a path that does not exist.
func parseURIList(data []byte) []string {
	var out []string
	for line := range strings.Lines(string(data)) {
		line = strings.TrimRight(line, "\r\n")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if path, ok := pathFromURI(line); ok {
			out = append(out, path)
		}
	}
	return out
}
