package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// SignAPK aligns and signs a package that is already built, with no config
// around it. It is what [Sign] does, and what "antuiapk sign" is: an APK
// that came from somewhere else — an unsigned one, or one to be re-signed
// with a release key — has no Go package behind it to describe.
func SignAPK(tc *sdk.Toolchain, params SignAPKParams) error {
	k := params.Keystore
	minSDK := params.MinSDK
	in := params.In
	out := params.Out
	if minSDK == 0 {
		minSDK = sdk.MinSDK
	}
	// The zero value is the debug key, which is what [Config.check] would
	// have done for a build. Signing on its own goes through neither.
	k.fill()
	if err := k.Ensure(tc); err != nil {
		return err
	}
	aligned := in + ".aligned"
	if _, err := run(tc.BuildTools.Zipalign(), "-P", "16", "-f", "4", in, aligned); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(mustAbs(out)), 0o755); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	_, err := run(tc.BuildTools.Apksigner(), "sign",
		"--ks", k.Path,
		"--ks-key-alias", k.Alias,
		"--ks-pass", "pass:"+k.StorePass,
		"--key-pass", "pass:"+k.KeyPass,
		"--min-sdk-version", fmt.Sprint(minSDK),
		// v4 is off. It is a detached signature in a second file beside the
		// APK, used only to make "adb install" incremental, and leaving a
		// stray .idsig next to whatever the caller asked for is a worse
		// surprise than a slower install.
		"--v4-signing-enabled", "false",
		"--out", out,
		aligned,
	)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}
