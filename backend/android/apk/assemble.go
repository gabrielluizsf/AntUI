package apk

import (
	"archive/zip"
	"fmt"
	"os"
	"path"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// Assemble copies the base APK and adds the libraries — and the Java shim,
// when the config asks for one — to it.
//
// The libraries go in **stored, not deflated**. That is not an optimisation:
// the manifest says extractNativeLibs="false", so the loader maps the .so
// straight out of the APK, and a compressed entry cannot be mapped. The
// alignment that mapping also needs is [Sign]'s job.
//
// Everything already in the base APK is copied byte for byte, keeping
// whatever compression aapt2 chose — except resources.arsc, which since API
// 30 must also be stored, and is decompressed here if it is not.
func Assemble(base string, libs map[sdk.ABI]string, libName, out string, cfg *Config) error {
	zr, err := zip.OpenReader(base)
	if err != nil {
		return fmt.Errorf("apk: reading %s: %w", base, err)
	}
	defer zr.Close()

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)

	for _, e := range zr.File {
		if e.Name == "resources.arsc" && e.Method != zip.Store {
			if err := storeEntry(zw, e); err != nil {
				return err
			}
			continue
		}
		if err := copyRaw(zw, e); err != nil {
			return err
		}
	}
	for abi, lib := range libs {
		name := path.Join("lib", string(abi), "lib"+libName+".so")
		if err := addStored(zw, name, lib); err != nil {
			return err
		}
	}
	// The shim, when the manifest says there is one. It is deflated rather
	// than stored: a dex is read by the runtime rather than mapped, so
	// nothing needs it uncompressed, and it compresses well.
	if cfg != nil && cfg.Shim != nil && *cfg.Shim {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: "classes.dex", Method: zip.Deflate})
		if err != nil {
			return fmt.Errorf("apk: classes.dex: %w", err)
		}
		if _, err := w.Write(shim.Dex); err != nil {
			return fmt.Errorf("apk: classes.dex: %w", err)
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return f.Close()
}
