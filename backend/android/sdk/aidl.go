package sdk

import (
	"path/filepath"
)

// Aidl is the tool that compiles an interface definition. Nothing here
// uses it; it is listed because a build-tools install without it is a
// broken one.
func (b BuildTools) Aidl() string { return filepath.Join(b.Dir, exe("aidl")) }
