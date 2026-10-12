package sdk

import (
	"path/filepath"
)

// Jar is the android.jar the Java shim compiles against.
func (p Platform) Jar() string { return filepath.Join(p.Dir, "android.jar") }
