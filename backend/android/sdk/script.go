package sdk

import (
	"runtime"
)

// script adds the extension a *wrapper* has: the SDK ships d8 and apksigner
// as shell scripts everywhere and as .bat files on Windows, which is not the
// same rule as a compiled tool like aapt2.
func script(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".bat"
	}
	return name
}
