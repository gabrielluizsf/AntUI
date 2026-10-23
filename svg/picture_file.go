package svg

import (
	"fmt"
	"net/url"
)

// pictureFile reads a `file:` address from this machine: the path it names,
// with the local host or none allowed and any other host refused, since a path
// on somebody else's machine is not one this can open. The path is taken as
// written after the scheme — `file:///a.png` and `file:a.png` both mean the
// file they name, absolute and relative respectively.
func pictureFile(target string) ([]byte, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != "localhost" {
		return nil, fmt.Errorf("the file address names the host %q, which is not this machine", u.Host)
	}
	path := filePathFromURL(u)
	if path == "" {
		// `file:a.png` carries no authority, and the path lands where Go keeps
		// an address that is only an opaque string.
		path = u.Opaque
	}
	return readLocalPicture(path)
}
