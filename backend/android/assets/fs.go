//go:build android

package assets

import "io/fs"

// FS is everything packed into the APK, as an ordinary [io/fs.FS].
//
// That is the whole point of this package. A game that reads its levels with
// os.ReadFile on a desktop reads them with fs.ReadFile here, and the same
// code serves both — which is the difference between porting a program and
// rewriting it.
//
//	levels, err := fs.ReadFile(assets.FS(), "levels/one.json")
//	err = fs.WalkDir(assets.FS(), "sounds", func(p string, d fs.DirEntry, err error) error { ... })
func FS() fs.FS { return assetFS{} }
