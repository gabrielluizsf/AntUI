//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/looper.h>
*/
import "C"

// Looper is Android's event loop for one thread. It is thread-local: the
// looper a thread prepares is the only one it can poll, and a file
// descriptor attached to it is only ever delivered there.
//
// That is why the app's goroutine is locked to an OS thread. Go is free to
// move a goroutine between threads, and a goroutine that moved after
// attaching the input queue would simply stop receiving input, with nothing
// anywhere to say why.
type Looper struct {
	ptr *C.ALooper
}
