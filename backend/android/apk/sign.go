package apk

import "github.com/gabrielluizsf/antui/backend/android/sdk"

// Sign aligns the package and signs it.
//
// The order is not free: zipalign moves entries, so signing first would
// invalidate the signature. apksigner knows this and refuses to sign
// something it would then have to be run over again.
//
// The alignment asked for is 16 KB, not the 4 KB that was standard for
// years. Android 15 introduced devices with 16 KB memory pages, and the
// Play Store requires an app to work on them; a library aligned to 4 KB
// cannot be mapped on such a device at all.
func Sign(tc *sdk.Toolchain, cfg *Config, paths Paths) error {
	if err := cfg.check(); err != nil {
		return err
	}
	return SignAPK(tc, SignAPKParams{Keystore: cfg.Keystore, MinSDK: cfg.MinSDK, In: paths.In, Out: paths.Out})
}
