//go:build android

package picker

import "os"

// Writer opens a chosen file for writing, replacing what is in it.
func Writer(uri string) (*os.File, error) { return openMode(uri, "wt") }
