//go:build android

package gallery

import "errors"

// Write puts bytes in.
func (w *Writer) Write(p []byte) (int, error) {
	if w.closed {
		return 0, errors.New("gallery: the entry is closed")
	}
	return w.file.Write(p)
}
