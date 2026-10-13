package sdk

import "path/filepath"

// D8 is the tool that turns Java class files into a dex.
func (b BuildTools) D8() string { return filepath.Join(b.Dir, script("d8")) }
