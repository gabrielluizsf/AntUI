package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

// compiledIcons returns an aapt2-compiled archive of the launcher icons,
// from the cache if it is there and freshly made if it is not.
//
// A cache that cannot be reached is not an error: it does the work and says
// nothing. Being unable to write to a cache directory is a reason to be
// slower and not a reason to fail a build.
func compiledIcons(tc *sdk.Toolchain, work, icon string, bg canvas.Color) (string, error) {
	key, err := iconKey(tc, icon, bg)
	cached := ""
	if err == nil {
		if dir, err := cacheDir(); err == nil {
			cached = filepath.Join(dir, key+".zip")
			if _, err := os.Stat(cached); err == nil {
				return cached, nil
			}
		}
	}

	res := filepath.Join(work, "res")
	if err := WriteIcons(res, icon, bg); err != nil {
		return "", err
	}
	out := filepath.Join(work, "resources.zip")
	if _, err := run(tc.BuildTools.Aapt2(), "compile", "--dir", res, "-o", out); err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	if cached != "" {
		// Written beside and renamed, so that two builds at once cannot
		// leave a half-written archive for a third to read.
		save(out, cached)
	}
	return out, nil
}
