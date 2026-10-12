//go:build android

package assets

// ReadFile is the whole of one asset.
func ReadFile(name string) ([]byte, error) { return assetFS{}.ReadFile(name) }
