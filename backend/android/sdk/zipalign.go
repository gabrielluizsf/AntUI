package sdk

import (
	"path/filepath"
)

// Zipalign is the tool that aligns a package so the loader can map what
// is inside it.
func (b BuildTools) Zipalign() string { return filepath.Join(b.Dir, exe("zipalign")) }
