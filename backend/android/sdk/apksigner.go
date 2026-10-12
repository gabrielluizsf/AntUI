package sdk

import (
	"path/filepath"
)

// Apksigner is the tool that signs an APK, and the only one that can
// write the v2 and v3 signatures a modern Android insists on.
func (b BuildTools) Apksigner() string { return filepath.Join(b.Dir, script("apksigner")) }
