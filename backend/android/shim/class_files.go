package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// classFiles lists what javac produced, in a fixed order.
func classFiles(classes string) ([]string, error) {
	var compiled []string
	err := filepath.WalkDir(classes, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(p, ".class") {
			compiled = append(compiled, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	if len(compiled) == 0 {
		return nil, fmt.Errorf("shim: javac produced no class files")
	}
	sort.Strings(compiled)
	return compiled, nil
}
