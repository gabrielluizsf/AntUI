package apk

import (
	"fmt"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// SignBundle signs with jarsigner.
//
// Not apksigner, which refuses a bundle: the signature schemes it writes go
// in an APK's own structures, and a bundle is signed the old way, as a JAR.
// What the user installs is signed later — by the store, with the key it
// holds — so this signature only proves who uploaded it.
func SignBundle(tc *sdk.Toolchain, cfg *Config, paths Paths) error {
	in := paths.In
	out := paths.Out
	if err := cfg.Keystore.Ensure(tc); err != nil {
		return err
	}
	k := cfg.Keystore
	jarsigner := filepath.Join(tc.JDK.Dir, "bin", "jarsigner")
	_, err := run(jarsigner,
		"-keystore", k.Path,
		"-storepass", k.StorePass,
		"-keypass", k.KeyPass,
		"-signedjar", out,
		"-digestalg", "SHA-256",
		"-sigalg", "SHA256withRSA",
		in, k.Alias,
	)
	if err != nil {
		return fmt.Errorf("apk: signing the bundle: %w", err)
	}
	return nil
}
