//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/looper.h>
*/
import "C"

// PrepareLooper makes this thread's looper, or returns the one it already
// has. It must be called from the thread that will poll it.
//
// The looper is prepared to allow polling without callbacks, which is what
// lets [Looper.Poll] hand an identifier back rather than calling into Go
// from C on a thread Go may not be expecting.
func PrepareLooper() *Looper {
	p := C.ALooper_prepare(C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	if p == nil {
		return nil
	}
	return &Looper{ptr: p}
}
