package apk

import (
	"archive/zip"
	"fmt"
)

func add(zw *zip.Writer, name string, body []byte) error {
	w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
	if err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	return nil
}
