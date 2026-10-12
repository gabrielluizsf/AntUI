//go:build android

package assets

import (
	"io/fs"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

func (f assetFS) ReadFile(name string) ([]byte, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrInvalid}
	}
	m, err := manager()
	if err != nil {
		return nil, &fs.PathError{Op: "read", Path: name, Err: err}
	}
	a, err := m.Open(inner(name), ndk.AssetBuffer)
	if err != nil {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
	}
	defer a.Close()
	return a.ReadAll()
}
