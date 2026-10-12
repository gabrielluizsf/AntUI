//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/looper.h>
*/
import "C"

// Wake makes a Poll that is waiting return straight away. It is safe from
// any thread, and is the only thing here that is.
func (l *Looper) Wake() { C.ALooper_wake(l.ptr) }
