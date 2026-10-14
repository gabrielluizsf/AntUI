package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Create makes a keystore with a fresh RSA key. days is how long the
// certificate is valid: the Play Store wants a key that outlives the app, so
// a real one is made for 10000 days — a shorter one cannot be extended, and
// an expired one cannot publish an update.
func Create(tc *sdk.Toolchain, params CreateParams) error {
	if tc.JDK.Dir == "" {
		return fmt.Errorf("apk: making a keystore needs a JDK, and none was found")
	}
	if err := os.MkdirAll(filepath.Dir(params.Keystore.Path), 0o700); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	keytool := filepath.Join(tc.JDK.Dir, "bin", "keytool")
	_, err := run(keytool,
		"-genkeypair",
		"-keystore", params.Keystore.Path,
		"-alias", params.Keystore.Alias,
		"-storepass", params.Keystore.StorePass,
		"-keypass", params.Keystore.KeyPass,
		"-keyalg", "RSA",
		"-keysize", "2048",
		"-validity", fmt.Sprint(params.Days),
		"-dname", params.Name,
		// PKCS12 is the default and the only format still current; saying so
		// keeps keytool from printing a migration warning on every build.
		"-storetype", "PKCS12",
	)
	if err != nil {
		return fmt.Errorf("apk: creating %s: %w", params.Keystore.Path, err)
	}
	return nil
}
