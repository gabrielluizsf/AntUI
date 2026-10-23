package svg

import (
	"fmt"
)

// fetchPicture is the bytes behind an address that has already been resolved
// against the base the drawing sits at: the scheme it names says how, through
// the table of schemes, and the scheme it does not name — an address with no
// scheme at all, the file a drawing beside it names while it is being written —
// is read as the file on this machine. A scheme the table does not hold is said
// out loud rather than read as a file name.
func fetchPicture(target string) ([]byte, error) {
	if scheme := hrefScheme(target); scheme != "" {
		fetch, ok := pictureSchemes[scheme]
		if !ok {
			return nil, fmt.Errorf("the scheme %q is not one a picture is read from", scheme+":")
		}
		return fetch(target)
	}
	return readLocalPicture(target)
}
