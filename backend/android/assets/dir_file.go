//go:build android

package assets

import (
	"io/fs"
)

// dirFile is a directory opened as a file, which fs.WalkDir does.
type dirFile struct {
	name    string
	entries []fs.DirEntry
	at      int
}
