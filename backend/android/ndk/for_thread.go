//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/looper.h>
*/
import "C"

// ForThread is the looper this thread already has, or nil.
func ForThread() *Looper {
	p := C.ALooper_forThread()
	if p == nil {
		return nil
	}
	return &Looper{ptr: p}
}
