package apk

import (
	"archive/zip"
	"fmt"
	"io"
)

// storeEntry rewrites a compressed entry as a stored one.
func storeEntry(zw *zip.Writer, e *zip.File) error {
	rc, err := e.Open()
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	defer rc.Close()
	w, err := zw.CreateHeader(&zip.FileHeader{Name: e.Name, Method: zip.Store})
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	return nil
}
