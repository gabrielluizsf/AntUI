package apk

import (
	"fmt"
	"os"
)

// warnPlayProblems says what the store would say, said now. A bundle is
// built to be uploaded, so it is worth hearing before the two minutes of
// building rather than after — but it is said and not enforced: a bundle
// built to be tried on an emulator wants x86_64, which no phone has and no
// store accepts, and that is a perfectly good reason to build one.
func warnPlayProblems(cfg Config) {
	for _, p := range cfg.PlayProblems() {
		fmt.Fprintln(os.Stderr, "antuiapk: "+p.String())
	}
}
