//go:build android

package assets

// info is what Stat gives back. There are no times and no owners inside a
// package, so there is nothing to report but the name, the size and whether
// it is a directory.
type info struct {
	name string
	size int64
	dir  bool
}
