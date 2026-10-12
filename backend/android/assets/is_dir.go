//go:build android

package assets

func (i info) IsDir() bool { return i.dir }
