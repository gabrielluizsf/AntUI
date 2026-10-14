package apk

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// AssembleBundle rearranges what aapt2 produced into the layout a bundle
// has, and adds the libraries, the dex and the configuration.
func AssembleBundle(cfg *Config, params AssembleBundleParams) error {
	proto := params.Proto
	libs := params.Libs
	out := params.Out
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
	if err := addBundleEntries(zw, &zr.Reader); err != nil {
		return err
	}
	if err := addBundleLibs(zw, libs, cfg); err != nil {
		return err
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
