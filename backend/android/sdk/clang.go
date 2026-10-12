package sdk

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

// Clang is the compiler that builds for one ABI against one API level. The
// NDK ships a wrapper per (target, level) pair rather than one compiler
// taking a flag, so the API level is baked into the name and a level with no
// wrapper is an error here rather than a confusing failure later.
func (n NDK) Clang(abi ABI, api int) (string, error) {
	if !abi.Valid() {
		return "", fmt.Errorf("sdk: %q is not an Android ABI", abi)
	}
	name := abi.Triple() + strconv.Itoa(api) + "-clang"
	if runtime.GOOS == "windows" {
		name += ".cmd"
	}
	path := filepath.Join(n.Bin(), name)
	if !exists(path) {
		return "", fmt.Errorf("sdk: NDK %s has no compiler for %s at API %d (%s)",
			n.Version, abi, api, name)
	}
	return path, nil
}
