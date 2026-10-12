package sdk

import (
	"path/filepath"
	"runtime"
)

func (t *Toolchain) cmdlineTool(name string) string {
	if runtime.GOOS == "windows" {
		name += ".bat"
	}
	return filepath.Join(t.Root, "cmdline-tools", "latest", "bin", name)
}
