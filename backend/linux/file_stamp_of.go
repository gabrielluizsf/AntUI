//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "os"

// fileStampOf reads a file's identity without opening it — whether it is
// there, how big it is, when it was last written — which is all it takes to
// notice that the desktop wrote a new answer into it. A file that is not
// there has a stamp of its own, so a desktop that stops writing settings
// is not asked about the ones it no longer keeps.
func fileStampOf(path string) fileStamp {
	info, err := os.Stat(path)
	if err != nil {
		return fileStamp{}
	}
	return fileStamp{size: info.Size(), mod: info.ModTime(), ok: true}
}
