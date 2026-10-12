package apk

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// playProblemsStore judges what the store checks before it will take an
// upload: the target API, the debug flag, the signing key and the id.
func playProblemsStore(cfg Config, ps []Problem) []Problem {
	if cfg.TargetSDK < sdk.PlayTarget {
		ps = append(ps, Problem{true,
			fmt.Sprintf("targets API %d, and the store requires %d",
				cfg.TargetSDK, sdk.PlayTarget),
			"the floor rises every August and applies to updates as well as " +
				"new apps"})
	}
	if cfg.Debuggable {
		ps = append(ps, Problem{true,
			"the app is debuggable",
			"android:debuggable lets anything on the device read the app's " +
				"data and attach to it; drop -debug"})
	}
	if cfg.Keystore.Debug() {
		ps = append(ps, Problem{true,
			"signed with the shared debug key",
			"every Android developer has that key and its password; make one " +
				`of your own with "antuiapk keygen"`})
	}
	if strings.HasPrefix(cfg.Package, "com.example.") || cfg.Package == "com.example" {
		ps = append(ps, Problem{true,
			fmt.Sprintf("the application id is %q", cfg.Package),
			"the store reserves com.example; and the id can never be changed " +
				"after the first release"})
	}
	return ps
}
