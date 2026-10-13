//go:build android

package assets

import "io/fs"

// ReadDir is what is directly inside a directory.
func ReadDir(name string) ([]fs.DirEntry, error) { return assetFS{}.ReadDir(name) }
