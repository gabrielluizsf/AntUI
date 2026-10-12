package apk

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// PlayProblems is what the store will say about this config, answered before
// anything is built rather than after it is uploaded.
//
// It reads the config and nothing else, so it is worth calling early: the
// whole point is to fail in a second instead of after a two-minute build and
// a day of waiting.
func (c Config) PlayProblems() []Problem {
	// A copy, so that asking the question does not change the answer: check
	// fills in the defaults, and the defaults are most of what is being
	// judged.
	cfg := c
	if err := cfg.check(); err != nil {
		return []Problem{{Fatal: true, Text: err.Error()}}
	}
	var ps []Problem

	if !carries64(cfg.ABIs) {
		ps = append(ps, Problem{true,
			fmt.Sprintf("no 64-bit ABI (%s)", abiList(cfg.ABIs)),
			"the store has refused 32-bit-only uploads since August 2019; " +
				"build for " + string(sdk.Arm64)})
	}
	if !slices.Contains(cfg.ABIs, sdk.Arm64) {
		ps = append(ps, Problem{false,
			"no arm64-v8a",
			"almost every phone sold since 2016 is arm64; without it the app " +
				"is not offered to them at all"})
	}
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
	if cfg.Icon == "" {
		ps = append(ps, Problem{false,
			"no icon",
			"the app appears under the platform's blank one, and the listing " +
				"needs a 512x512 one separately"})
	}
	if cfg.Strip != nil && !*cfg.Strip {
		ps = append(ps, Problem{false,
			"the library ships with its symbols in it",
			"it is about a third larger than it needs to be, and the store " +
				"asks for the symbols as a separate file anyway"})
	}
	if cfg.VersionCode <= 0 {
		ps = append(ps, Problem{true,
			"no version code",
			"it must be a positive number and must rise with every upload; a " +
				"code that has been used once can never be used again"})
	}
	return ps
}
