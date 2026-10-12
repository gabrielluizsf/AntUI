package sdk

import (
	"os"
	"path/filepath"
)

func findBuildTools(root string) (BuildTools, bool) {
	var found []BuildTools
	for _, name := range subdirs(filepath.Join(root, "build-tools")) {
		dir := filepath.Join(root, "build-tools", name)
		// A build-tools directory that lost its aapt2 is a half-finished
		// download, and picking it would fail much later and less clearly.
		if _, err := os.Stat(filepath.Join(dir, exe("aapt2"))); err != nil {
			continue
		}
		found = append(found, BuildTools{Dir: dir, Version: ParseVersion(name)})
	}
	return pick(found, func(b BuildTools) Version { return b.Version })
}
