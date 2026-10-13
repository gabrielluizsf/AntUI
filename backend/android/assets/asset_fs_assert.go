//go:build android

package assets

import "io/fs"

var (
	_ fs.FS         = assetFS{}
	_ fs.ReadDirFS  = assetFS{}
	_ fs.ReadFileFS = assetFS{}
)
