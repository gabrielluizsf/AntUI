package apk

import (
	"fmt"
	"slices"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// playProblemsABIs judges the ABI list: the store refuses 32-bit-only
// uploads, and arm64 is what almost every phone needs.
func playProblemsABIs(cfg Config, ps []Problem) []Problem {
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
	return ps
}
