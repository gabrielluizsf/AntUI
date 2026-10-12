package shim

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// writeSources writes the embedded Java files into work and returns their
// paths, in a fixed order so the class files reach d8 the same way twice.
func writeSources(work string) ([]string, error) {
	names, err := fs.Glob(Sources, "*.java")
	if err != nil || len(names) == 0 {
		return nil, fmt.Errorf("shim: no Java to compile")
	}
	sort.Strings(names)
	srcs := make([]string, 0, len(names))
	for _, name := range names {
		b, err := Sources.ReadFile(name)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(work, name)
		if err := os.WriteFile(path, b, 0o644); err != nil {
			return nil, err
		}
		srcs = append(srcs, path)
	}
	return srcs, nil
}
