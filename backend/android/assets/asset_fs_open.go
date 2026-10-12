//go:build android

package assets

import (
	"io/fs"
	"path"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

func (f assetFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	m, err := manager()
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	// A directory first: opening one as a file gives an error from the
	// platform that says nothing about it being a directory.
	if entries, err := list(name); err == nil && len(entries) > 0 {
		return &dirFile{name: path.Base(name), entries: entries}, nil
	}
	// Buffered, so that an uncompressed asset can be handed over as the
	// mapping it already is rather than copied.
	a, err := m.Open(inner(name), ndk.AssetBuffer)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &file{name: path.Base(name), asset: a}, nil
}
