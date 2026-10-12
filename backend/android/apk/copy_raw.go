package apk

import (
	"archive/zip"
	"fmt"
	"io"
)

// copyRaw moves an entry across without recompressing it.
func copyRaw(zw *zip.Writer, e *zip.File) error {
	rc, err := e.OpenRaw()
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	h := e.FileHeader
	w, err := zw.CreateRaw(&h)
	if err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("apk: %s: %w", e.Name, err)
	}
	return nil
}
