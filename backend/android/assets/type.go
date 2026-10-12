//go:build android

package assets

import (
	"io/fs"
)

func (e entry) Type() fs.FileMode { return e.info.Mode().Type() }
