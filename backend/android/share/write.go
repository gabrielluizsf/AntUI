//go:build android

package share

// Write makes a file with these bytes in it and gives back its address.
func Write(name string, data []byte) (string, error) {
	f, uri, err := Create(name)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return uri, nil
}
