package sdk

import (
	"os"
	"path/filepath"
)

// subdirs lists a directory's subdirectories, sorted, and gives back nothing
// at all when the directory does not exist.
func subdirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		// A symlinked SDK directory is common enough to be worth following.
		if e.IsDir() {
			names = append(names, e.Name())
			continue
		}
		if e.Type()&os.ModeSymlink != 0 {
			if st, err := os.Stat(filepath.Join(dir, e.Name())); err == nil && st.IsDir() {
				names = append(names, e.Name())
			}
		}
	}
	return names
}
