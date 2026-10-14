//go:build android

package gallery

// Save writes bytes into the device's library in one call, and gives back
// the content:// address.
func Save(entry Entry, data []byte) (string, error) {
	w, err := Create(entry)
	if err != nil {
		return "", err
	}
	if _, err := w.Write(data); err != nil {
		w.abandon()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return w.URI(), nil
}
