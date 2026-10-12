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
	// Resources, when there are any. aapt2 works in two steps and they are
	// not interchangeable: compile turns a directory of files into a flat
	// archive of compiled ones, and link takes those archives and the
	// manifest and makes the package. Handing link a directory of PNGs does
	// nothing at all.
	var compiled string
	if cfg.Icon != "" {
		var err error
		// The icons and their compiled form depend on nothing but the
		// picture, so this is usually a cache read — see [compiledIcons].
		compiled, err = compiledIcons(tc, dir, cfg.Icon, cfg.IconBackground)
		if err != nil {
			return err
		}
	}

	args := []string{"link",
		"-I", tc.Platform.Jar(),
		"--manifest", manifest,
		"-o", out,
		"--min-sdk-version", fmt.Sprint(cfg.MinSDK),
		"--target-sdk-version", fmt.Sprint(cfg.TargetSDK),
		"--version-code", fmt.Sprint(cfg.VersionCode),
		"--version-name", cfg.VersionName,
	}
	if cfg.Assets != "" {
		args = append(args, "-A", cfg.Assets)
		// Anything already compressed goes in stored, so that it can be
		// mapped rather than inflated on every read.
		for _, ext := range cfg.Store {
			args = append(args, "-0", ext)
		}
	}
	// The compiled resources go on the end as a positional argument, not
	// behind a flag. -R exists and means something else: it is for *overlay*
	// resources, and passing the ordinary ones to it fails with "failed
	// parsing overlays", which does not suggest what to do about it.
	if proto {
		// The manifest and the resource table come out as protocol buffers
		// rather than binary XML and an arsc, which is what a bundle holds.
		args = append(args, "--proto-format")
	}
	if compiled != "" {
		args = append(args, compiled)
	}
	if _, err := run(tc.BuildTools.Aapt2(), args...); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}
