//go:build android

package assets

import "io/fs"

// Open is fs.FS.Open on [FS].
func Open(name string) (fs.File, error) { return assetFS{}.Open(name) }
