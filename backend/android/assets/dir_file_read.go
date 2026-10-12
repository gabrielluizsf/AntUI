//go:build android

package assets

import (
	"errors"
	"io/fs"
)

func (d *dirFile) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.name, Err: errors.New("is a directory")}
}
