//go:build android

package picker

import "os"

// Open reads a file the user chose.
//
// The address is a content:// one, not a path, and it may have no file
// behind it at all — the document could be in a cloud account, being
// streamed as it is read. So this is the only way to get at one, and the
// file it gives back is a descriptor rather than a name on disk.
func Open(uri string) (*os.File, error) { return openMode(uri, "r") }
