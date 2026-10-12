//go:build android

package ndk

// Size is how many bytes the asset is.
func (a *Asset) Size() int64 { return a.size }
