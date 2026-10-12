//go:build android

package assets

import (
	"io/fs"
)

func (i info) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
