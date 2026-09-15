//go:build android

package ndk

import (
	"bufio"
	"os"
	"sync"
	"syscall"
)

var redirectOnce sync.Once

// RedirectStdio points file descriptors 1 and 2 at logcat.
//
// This matters more on Android than it looks. An app has no terminal, so
// everything a Go program writes to stdout or stderr — fmt.Println, the log
// package's default, and every panic traceback the runtime prints — goes
// nowhere at all. After this call it goes to logcat under [Tag], and
// "adb logcat -s antui" is a usable console.
//
// The one thing it cannot promise is the tail of a crash: the runtime writes
// the traceback and then kills the process, and the goroutine reading the
// pipe may not be scheduled again before that happens. Most of a traceback
// arrives; the last lines sometimes do not.
//
// Calling it more than once does nothing.
func RedirectStdio() {
	redirectOnce.Do(func() {
		r, w, err := os.Pipe()
		if err != nil {
			Errorf("cannot redirect stdio: %v", err)
			return
		}
		// dup3 rather than dup2: arm64 Linux has no dup2 syscall at all, and
		// Go's syscall package reflects that.
		fd := int(w.Fd())
		if err := syscall.Dup3(fd, 1, 0); err != nil {
			Errorf("cannot redirect stdout: %v", err)
			return
		}
		if err := syscall.Dup3(fd, 2, 0); err != nil {
			Errorf("cannot redirect stderr: %v", err)
			return
		}
		go pump(r)
	})
}

// pump reads the pipe a line at a time and files each line under the tag.
// A line longer than the buffer is logged in pieces rather than dropped,
// because the piece that gets dropped is always the interesting one.
func pump(r *os.File) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 8*1024), 512*1024)
	for sc.Scan() {
		Log(Info, Tag, sc.Text())
	}
	if err := sc.Err(); err != nil {
		Log(Error, Tag, "stdio: "+err.Error())
	}
}
