//go:build android

package assets

func (e entry) IsDir() bool { return e.info.dir }
