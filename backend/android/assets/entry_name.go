//go:build android

package assets

func (e entry) Name() string { return e.info.name }
