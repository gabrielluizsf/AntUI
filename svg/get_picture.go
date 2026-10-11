package svg

import (
	"fmt"
	"io"
	"net/http"
)

// getPicture is one address read whole, or the reason it could not be: a
// server that did not answer with the picture is said by the answer it did
// give, and a body past the size this will take is refused before it is in
// memory rather than after.
func getPicture(target string) ([]byte, error) {
	resp, err := pictureHTTP.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("it answered %s", resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxPictureBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxPictureBytes {
		return nil, fmt.Errorf("it is bigger than %d bytes", maxPictureBytes)
	}
	return data, nil
}
