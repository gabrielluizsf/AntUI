//go:build android

package share

import (
	"os"
	"path/filepath"
)

// Clear empties the directory. Worth doing when the app is done handing
// things out: the copies are of no further use and the system will otherwise
// keep them until it needs the space.
func Clear() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		os.RemoveAll(filepath.Join(dir, e.Name()))
	}
	return nil
}
