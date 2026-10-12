package shim

import (
	"fmt"
	"os/exec"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// runD8 turns the class files into the dex. d8 is deterministic given the
// same inputs in the same order, which is what makes the committed dex
// checkable.
func runD8(tc *sdk.Toolchain, work string, compiled []string) error {
	d8 := append([]string{"--lib", tc.Platform.Jar(), "--output", work}, compiled...)
	if out, err := exec.Command(tc.BuildTools.D8(), d8...).CombinedOutput(); err != nil {
		return fmt.Errorf("shim: d8: %w\n%s", err, out)
	}
	return nil
}
