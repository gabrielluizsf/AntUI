//go:build android

package assets

func (f *file) Close() error { return f.asset.Close() }
