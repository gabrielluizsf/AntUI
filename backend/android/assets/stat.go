//go:build android

package assets

import "io/fs"

func (f *file) Stat() (fs.FileInfo, error) {
	return info{name: f.name, size: f.asset.Size()}, nil
}
