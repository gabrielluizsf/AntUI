package apk

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

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
