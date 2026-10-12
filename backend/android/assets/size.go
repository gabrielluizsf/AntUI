//go:build android

package assets

func (i info) Size() int64 { return i.size }
