package sdk

import (
	"path/filepath"
)

// The four tools an APK is built with. Two of them, d8 and apksigner, are
// wrapper scripts rather than programs, which only shows on Windows.
func (b BuildTools) Aapt2() string { return filepath.Join(b.Dir, exe("aapt2")) }
