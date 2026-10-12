package apk

import (
	"archive/zip"
	"fmt"
	"os"
	"path"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// addBundleLibs puts the shared library for each ABI under the module the
// bundle keeps them in.
func addBundleLibs(zw *zip.Writer, libs map[sdk.ABI]string, cfg *Config) error {
	for abi, lib := range libs {
		body, err := os.ReadFile(lib)
		if err != nil {
			return fmt.Errorf("apk: %w", err)
		}
		name := path.Join("base", "lib", string(abi), "lib"+cfg.Lib+".so")
		if err := add(zw, name, body); err != nil {
			return err
		}
	}
	return nil
}
