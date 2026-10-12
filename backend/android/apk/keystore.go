package apk

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
