//go:build android

package ndk

/*
#cgo LDFLAGS: -llog
#include <android/log.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// Priority is how loud a log line is. These are the levels logcat filters on.
type Priority int32

// The priorities, as logcat spells them.
const (
	Verbose Priority = 2
	Debug   Priority = 3
	Info    Priority = 4
	Warn    Priority = 5
	Error   Priority = 6
	Fatal   Priority = 7
)

// Tag is what every line from this library is filed under, unless a caller
// says otherwise. Keeping it short and constant is what makes
// "adb logcat -s antui" useful.
const Tag = "antui"

// Log writes one line to the system log. A line longer than the kernel's
// buffer is truncated by the platform, not here — splitting it would
// interleave with other processes' lines and read worse than a cut one.
func Log(prio Priority, tag, text string) {
	ctag := C.CString(tag)
	ctext := C.CString(text)
	C.__android_log_write(C.int(prio), ctag, ctext)
	C.free(unsafe.Pointer(ctag))
	C.free(unsafe.Pointer(ctext))
}

// Logf is Log with formatting. It is deliberately not variadic-over-Log so
// that the common case of a plain string never formats.
func Logf(prio Priority, format string, args ...any) {
	Log(prio, Tag, sprintf(format, args...))
}

// Infof, Warnf and Errorf are the three levels this library actually uses.
func Infof(format string, args ...any) { Logf(Info, format, args...) }

// Warnf writes a line to logcat at the warning level.
func Warnf(format string, args ...any) { Logf(Warn, format, args...) }

// Errorf writes a line to logcat at the error level, which is the one a
// crash report keeps.
func Errorf(format string, args ...any) { Logf(Error, format, args...) }
