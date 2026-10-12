package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// VerifyBundle takes a bundle apart the way the store will, and puts the
// result on a device.
//
// This is the only check that means anything about a bundle. Everything
// before it says the file is shaped right; this says the store can turn it
// into an app that installs and runs — which is a different question, and
// the one that matters at three in the morning before a release.
//
// serial names a device, or is empty for the only one attached.
func VerifyBundle(tc *sdk.Toolchain, cfg *Config, bundle, serial string) error {
	if err := cfg.check(); err != nil {
		return err
	}
	if tc.Bundletool == "" {
		return ErrNoBundletool
	}
	if tc.JDK.Dir == "" {
		return fmt.Errorf("apk: bundletool is a jar and needs a JDK to run")
	}

	work, err := os.MkdirTemp("", "antuiverify-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	apks := filepath.Join(work, "out.apks")

	k := cfg.Keystore
	args := []string{"-jar", tc.Bundletool, "build-apks",
		"--bundle=" + bundle,
		"--output=" + apks,
		"--ks=" + k.Path,
		"--ks-key-alias=" + k.Alias,
		"--ks-pass=pass:" + k.StorePass,
		"--key-pass=pass:" + k.KeyPass,
		// Only the splits this device needs, which is what the store does
		// and is the whole point of the format.
		"--connected-device",
	}
	if serial != "" {
		args = append(args, "--device-id="+serial)
	}
	java := filepath.Join(tc.JDK.Dir, "bin", "java")
	if _, err := run(java, args...); err != nil {
		return fmt.Errorf("apk: bundletool could not build APKs from the bundle: %w", err)
	}

	install := []string{"-jar", tc.Bundletool, "install-apks", "--apks=" + apks}
	if serial != "" {
		install = append(install, "--device-id="+serial)
	}
	if _, err := run(java, install...); err != nil {
		return fmt.Errorf("apk: the APKs the bundle produced would not install: %w", err)
	}
	return nil
}
