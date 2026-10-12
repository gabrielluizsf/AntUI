package apk

import (
	"archive/zip"
	"fmt"
	"io"
)

// addBundleEntries copies what aapt2 produced into the bundle layout,
// dropping whatever is not part of it.
func addBundleEntries(zw *zip.Writer, zr *zip.Reader) error {
	for _, e := range zr.File {
		name := bundlePath(e.Name)
		if name == "" {
			continue
		}
		rc, err := e.Open()
		if err != nil {
			return fmt.Errorf("apk: %s: %w", e.Name, err)
		}
		body, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("apk: %s: %w", e.Name, err)
		}
		if err := add(zw, name, body); err != nil {
			return err
		}
	}
	return nil
}
