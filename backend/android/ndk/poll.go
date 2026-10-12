//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/looper.h>
*/
import "C"

// Poll waits for something to happen, for at most timeout milliseconds — 0
// to look and return, and a negative number to wait forever.
//
// It returns the identifier the source was attached with, or one of the
// Poll* constants above. Note that everything except an identifier is
// negative, so a caller that only cares about its own source can test for
// that identifier and ignore the rest.
func (l *Looper) Poll(timeoutMillis int) int {
	return int(C.ALooper_pollOnce(C.int(timeoutMillis), nil, nil, nil))
}
