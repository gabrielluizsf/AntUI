package apk

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func addStored(zw *zip.Writer, name, from string) error {
	src, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer src.Close()
	w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
	if err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	return nil
}
