package apk

import (
	"archive/zip"
	"fmt"
	"os"
	"path"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Assemble copies the base APK and adds the libraries — and the Java shim,
// when the config asks for one — to it.
//
// The libraries go in **stored, not deflated**. That is not an optimisation:
// the manifest says extractNativeLibs="false", so the loader maps the .so
// straight out of the APK, and a compressed entry cannot be mapped. The
// alignment that mapping also needs is [Sign]'s job.
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

	if err := copyBase(zw, &zr.Reader); err != nil {
		return err
	}
	for abi, lib := range libs {
		name := path.Join("lib", string(abi), "lib"+libName+".so")
		if err := addStored(zw, name, lib); err != nil {
			return err
		}
	}
	if err := addShim(zw, cfg); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return f.Close()
}
