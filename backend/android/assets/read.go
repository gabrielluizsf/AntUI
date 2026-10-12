//go:build android

package assets

func (f *file) Read(p []byte) (int, error) { return f.asset.Read(p) }
