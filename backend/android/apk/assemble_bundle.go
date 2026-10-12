package apk

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// AssembleBundle rearranges what aapt2 produced into the layout a bundle
// has, and adds the libraries, the dex and the configuration.
func AssembleBundle(proto string, libs map[sdk.ABI]string, cfg *Config, out string) error {
	zr, err := zip.OpenReader(proto)
	if err != nil {
		return fmt.Errorf("apk: reading %s: %w", proto, err)
	}
	defer zr.Close()

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)

	// BundleConfig first, at the root and outside every module.
	if err := add(zw, "BundleConfig.pb", bundleConfig()); err != nil {
		return err
	}

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
	if cfg.Shim != nil && *cfg.Shim {
		if err := add(zw, "base/dex/classes.dex", shim.Dex); err != nil {
			return err
		}
	}
	// The symbols ride along, under the path the store looks for them at.
	// They are metadata and not part of any module, so nothing about them
	// reaches a device: the bundle is where they are kept, not shipped.
	if err := addSymbols(zw, libs, cfg.Lib); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return f.Close()
}
