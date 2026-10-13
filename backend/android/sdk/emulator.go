package sdk

import "path/filepath"

// Emulator is the emulator, which is a separate download from the SDK.
func (t *Toolchain) Emulator() string { return filepath.Join(t.Root, "emulator", exe("emulator")) }
