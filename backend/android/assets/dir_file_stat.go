//go:build android

package assets

import "io/fs"

func (d *dirFile) Stat() (fs.FileInfo, error) { return info{name: d.name, dir: true}, nil }
