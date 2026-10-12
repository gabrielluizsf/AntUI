package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// link is Link and LinkProto, which differ by one flag and nothing else.
func link(tc *sdk.Toolchain, cfg *Config, out string, proto bool) error {
	if err := cfg.check(); err != nil {
		return err
	}
	dir := filepath.Dir(out)
	manifest := filepath.Join(dir, "AndroidManifest.xml")
	if err := os.WriteFile(manifest, []byte(cfg.Manifest()), 0o644); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	compiled := ""
	if cfg.Icon != "" {
		var err error
		// The icons and their compiled form depend on nothing but the
		// picture, so this is usually a cache read — see [compiledIcons].
		compiled, err = compiledIcons(tc, dir, cfg.Icon, cfg.IconBackground)
		if err != nil {
			return err
		}
	}
	if _, err := run(tc.BuildTools.Aapt2(), linkArgs(tc, cfg, manifest, out, compiled, proto)...); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}
