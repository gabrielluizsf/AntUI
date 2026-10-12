package apk

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
)

// addSymbol puts one symbol file in a zip under dir.
func addSymbol(z *zip.Writer, dir, lib, src string) error {
	name := path.Join(dir, "lib"+lib+".so"+symbolSuffix)
	w, err := z.Create(name)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer in.Close()
	if _, err := io.Copy(w, in); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}
