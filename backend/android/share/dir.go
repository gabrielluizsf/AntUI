//go:build android

package share

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/app"
)

// Dir is where files to be handed to other apps live.
//
// It is inside the cache, which means the system may empty it when the
// device runs short — and that is right: a file put here is a copy made to
// be given away, and if it goes, it goes.
func Dir() (string, error) {
	a := app.Current()
	if a == nil {
		return "", app.ErrNoUIThread
	}
	cache := a.Info().Cache
	if cache == "" {
		return "", errors.New("share: the app does not know its cache directory")
	}
	dir := filepath.Join(cache, dirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}
