package apk

import (
	"fmt"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// buildApksFromBundle turns a bundle into the splits one device needs, which
// is what the store does and is the whole point of the format.
func buildApksFromBundle(tc *sdk.Toolchain, cfg *Config, bundle, serial, apks string) error {
	k := cfg.Keystore
	args := []string{"-jar", tc.Bundletool, "build-apks",
		"--bundle=" + bundle,
		"--output=" + apks,
		"--ks=" + k.Path,
		"--ks-key-alias=" + k.Alias,
		"--ks-pass=pass:" + k.StorePass,
		"--key-pass=pass:" + k.KeyPass,
		"--connected-device",
	}
	if serial != "" {
		args = append(args, "--device-id="+serial)
	}
	java := filepath.Join(tc.JDK.Dir, "bin", "java")
	if _, err := run(java, args...); err != nil {
		return fmt.Errorf("apk: bundletool could not build APKs from the bundle: %w", err)
	}
	return nil
}
