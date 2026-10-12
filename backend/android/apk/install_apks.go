package apk

import (
	"fmt"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// installApks puts the APKs on the device, so that the check ends with the
// app actually installed rather than merely well-formed.
func installApks(tc *sdk.Toolchain, serial, apks string) error {
	java := filepath.Join(tc.JDK.Dir, "bin", "java")
	args := []string{"-jar", tc.Bundletool, "install-apks", "--apks=" + apks}
	if serial != "" {
		args = append(args, "--device-id="+serial)
	}
	if _, err := run(java, args...); err != nil {
		return fmt.Errorf("apk: the APKs the bundle produced would not install: %w", err)
	}
	return nil
}
