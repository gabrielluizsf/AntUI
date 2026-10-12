//go:build android

package assets

import (
	"io/fs"
)

func (e entry) Info() (fs.FileInfo, error) { return e.info, nil }
