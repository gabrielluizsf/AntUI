//go:build android

package assets

// Seek is here because a great many readers want it — an image decoder, a zip
// reader — and an uncompressed asset can do it.
func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.asset.Seek(offset, whence)
}
