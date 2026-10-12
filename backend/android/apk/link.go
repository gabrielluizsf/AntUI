package apk

import (
	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Link turns the manifest into a base APK with aapt2. The result holds the
// manifest in binary form and a resource table, and nothing else — no code
// and no libraries.
func Link(tc *sdk.Toolchain, cfg *Config, out string) error {
	return link(tc, cfg, out, false)
}
