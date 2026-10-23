package svg

import (
	"fmt"
	"io"
	"os"
)

// readLocalPicture reads a picture off this machine's disk, held to the same
// [maxPictureBytes] a server's answer is. Only a plain file is read: a path
// may name a device or a pipe as easily as a file, and `/dev/zero` read whole
// is the memory of the machine, `/dev/stdin` or a pipe nobody writes to a
// drawing that never finishes being read.
func readLocalPicture(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("it is not a plain file")
	}
	if st.Size() > maxPictureBytes {
		return nil, fmt.Errorf("it is bigger than %d bytes", maxPictureBytes)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxPictureBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxPictureBytes {
		return nil, fmt.Errorf("it is bigger than %d bytes", maxPictureBytes)
	}
	return data, nil
}
