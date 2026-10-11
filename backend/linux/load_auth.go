//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"strconv"
)

// loadAuth digs the MIT-MAGIC-COOKIE-1 for this display out of ~/.Xauthority.
// A missing or unreadable file is not an error: a server on a local socket
// often needs no cookie at all, and the handshake will say so if it does.
func loadAuth(displayNumber int) xAuth {
	path := authPath()
	if path == "" {
		return xAuth{}
	}
	f, err := os.Open(path)
	if err != nil {
		return xAuth{}
	}
	defer f.Close()
	return findAuthEntry(f, strconv.Itoa(displayNumber))
}
