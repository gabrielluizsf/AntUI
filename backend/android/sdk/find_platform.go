package sdk

import (
	"os"
	"path/filepath"
)

func findPlatform(root string) (Platform, bool) {
	var found []Platform
	for _, name := range subdirs(filepath.Join(root, "platforms")) {
		p, ok := parsePlatform(name)
		if !ok {
			continue
		}
		p.Dir = filepath.Join(root, "platforms", name)
		if _, err := os.Stat(p.Jar()); err != nil {
			continue
		}
		// An -ext platform is an older API level with later additions, so it
		// is never the right automatic choice; it can still be named by hand.
		if p.Ext != 0 {
			continue
		}
		found = append(found, p)
	}
	return pick(found, func(p Platform) Version {
		return Version{Parts: []int{p.API, p.Minor}}
	})
}
