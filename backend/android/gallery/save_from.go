//go:build android

package gallery

import "io"

// SaveFrom copies from a reader, for something too big to hold in memory.
func SaveFrom(entry Entry, r io.Reader) (string, error) {
	w, err := Create(entry)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(w, r); err != nil {
		w.abandon()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return w.URI(), nil
}
