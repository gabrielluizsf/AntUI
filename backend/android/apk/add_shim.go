package apk

import (
	"archive/zip"
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// addShim adds the compiled Java shim, when the manifest says there is one.
// It is deflated rather than stored: a dex is read by the runtime rather than
// mapped, so nothing needs it uncompressed, and it compresses well.
func addShim(zw *zip.Writer, cfg *Config) error {
	if cfg == nil || cfg.Shim == nil || !*cfg.Shim {
		return nil
	}
	w, err := zw.CreateHeader(&zip.FileHeader{Name: "classes.dex", Method: zip.Deflate})
	if err != nil {
		return fmt.Errorf("apk: classes.dex: %w", err)
	}
	if _, err := w.Write(shim.Dex); err != nil {
		return fmt.Errorf("apk: classes.dex: %w", err)
	}
	return nil
}
