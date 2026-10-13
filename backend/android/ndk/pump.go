//go:build android

package ndk

import (
	"bufio"
	"os"
)

const (
	logScannerInitialBytes = 8 * 1024
	logScannerMaxBytes     = 512 * 1024
)

// pump reads the pipe a line at a time and files each line under the tag.
// A line longer than the buffer is logged in pieces rather than dropped,
// because the piece that gets dropped is always the interesting one.
func pump(r *os.File) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, logScannerInitialBytes), logScannerMaxBytes)
	for sc.Scan() {
		Log(Info, Tag, sc.Text())
	}
	if err := sc.Err(); err != nil {
		Log(Error, Tag, "stdio: "+err.Error())
	}
}
