package sdk

import "path/filepath"

// Adb, and the rest of what lives outside build-tools.
func (t *Toolchain) Adb() string { return filepath.Join(t.Root, "platform-tools", exe("adb")) }
