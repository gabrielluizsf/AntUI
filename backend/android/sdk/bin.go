package sdk

import (
	"path/filepath"
)

// Bin is the NDK's toolchain directory, where clang and the llvm tools are.
func (n NDK) Bin() string {
	return filepath.Join(n.Dir, "toolchains", "llvm", "prebuilt", n.Host, "bin")
}
