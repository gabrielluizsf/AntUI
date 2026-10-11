package svg

import (
	"os"
	"strings"
)

// fetchPicture is the bytes behind an href: a `data:` URI unwrapped where it
// stands, an address fetched over the network — a protocol-relative one being
// the https it is meant to be, since there is no page here to take a scheme
// from — and anything else read as the file it names, which is what makes a
// picture beside the drawing work while it is being written.
func fetchPicture(href string) ([]byte, error) {
	if strings.HasPrefix(href, "data:") {
		return dataURI(href)
	}
	if target, ok := pictureHTTPURL(href); ok {
		return getPicture(target)
	}
	return os.ReadFile(href)
}
