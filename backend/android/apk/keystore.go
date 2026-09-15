package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Keystore is the key an APK is signed with. Android will not install an
// unsigned package, so every build signs — the question is only with what.
//
// The zero value is the debug keystore: the same well-known key, password
// and alias that every Android tool has used for fifteen years, kept at
// ~/.android/debug.keystore and created on first use. It is fine for
// installing on a device you own and useless for anything else, because
// everyone has it.
type Keystore struct {
	Path  string `json:"path,omitempty"`
	Alias string `json:"alias,omitempty"`
	// The passwords are deliberately not part of the written form: a
	// config file sits next to the source and gets committed, and a signing
	// key whose password is in the repository is not a signing key any more.
	// They come from the flags or from ANTUIAPK_STOREPASS and
	// ANTUIAPK_KEYPASS — see [Config.Load].
	StorePass string `json:"-"`
	KeyPass   string `json:"-"`
}

// The debug key, as the rest of the Android world spells it. These are not
// secrets and are not treated as any.
const (
	debugAlias = "androiddebugkey"
	debugPass  = "android"
	debugName  = "CN=Android Debug, O=Android, C=US"
)

// Debug reports whether this is the shared debug key rather than a real one.
func (k Keystore) Debug() bool { return k.Alias == debugAlias && k.StorePass == debugPass }

func (k *Keystore) fill() {
	if k.Path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			k.Path = filepath.Join(home, ".android", "debug.keystore")
		}
		if k.Alias == "" {
			k.Alias = debugAlias
		}
		if k.StorePass == "" {
			k.StorePass = debugPass
		}
	}
	if k.KeyPass == "" {
		k.KeyPass = k.StorePass
	}
}

// Ensure creates the keystore if it is the debug one and is not there yet.
// A real keystore is never created behind the caller's back: losing the key
// an app was published with means the app can never be updated again, so one
// is only ever made when asked for, by [Create].
func (k Keystore) Ensure(tc *sdk.Toolchain) error {
	if k.Path == "" {
		return fmt.Errorf("apk: no keystore")
	}
	if _, err := os.Stat(k.Path); err == nil {
		return nil
	}
	if !k.Debug() {
		return fmt.Errorf("apk: no keystore at %s; make one with antuiapk keygen", k.Path)
	}
	return Create(tc, k, debugName, 10000)
}

// Create makes a keystore with a fresh RSA key. days is how long the
// certificate is valid: the Play Store wants a key that outlives the app, so
// a real one is made for 10000 days — a shorter one cannot be extended, and
// an expired one cannot publish an update.
func Create(tc *sdk.Toolchain, k Keystore, name string, days int) error {
	if tc.JDK.Dir == "" {
		return fmt.Errorf("apk: making a keystore needs a JDK, and none was found")
	}
	if err := os.MkdirAll(filepath.Dir(k.Path), 0o700); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	keytool := filepath.Join(tc.JDK.Dir, "bin", "keytool")
	_, err := run(keytool,
		"-genkeypair",
		"-keystore", k.Path,
		"-alias", k.Alias,
		"-storepass", k.StorePass,
		"-keypass", k.KeyPass,
		"-keyalg", "RSA",
		"-keysize", "2048",
		"-validity", fmt.Sprint(days),
		"-dname", name,
		// PKCS12 is the default and the only format still current; saying so
		// keeps keytool from printing a migration warning on every build.
		"-storetype", "PKCS12",
	)
	if err != nil {
		return fmt.Errorf("apk: creating %s: %w", k.Path, err)
	}
	return nil
}
