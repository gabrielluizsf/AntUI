//go:build android

package assets

// assetFS satisfies fs.FS, fs.ReadDirFS and fs.ReadFileFS.
type assetFS struct{}
