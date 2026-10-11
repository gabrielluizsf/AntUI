package svg

import (
	"net/http"
	"time"
)

// pictureHTTP is what fetches a picture from an address, with a timeout: a
// library that reads a drawing synchronously cannot wait for a server for
// ever, and a picture that does not arrive in ten seconds is a picture that is
// not going to arrive in a way this drawing can use.
var pictureHTTP = &http.Client{Timeout: 10 * time.Second}
