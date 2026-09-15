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

// What Poll returns when it is not returning an identifier.
const (
	// PollWake means someone called Wake.
	PollWake = -1
	// PollCallback means the looper ran a callback itself; there is nothing
	// for the caller to do.
	PollCallback = -2
	// PollTimeout means the timeout ran out with nothing to report.
	PollTimeout = -3
	// PollError means the looper is broken.
	PollError = -4
)

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

// ForThread is the looper this thread already has, or nil.
func ForThread() *Looper {
	p := C.ALooper_forThread()
	if p == nil {
		return nil
	}
	return &Looper{ptr: p}
}

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

// Wake makes a Poll that is waiting return straight away. It is safe from
// any thread, and is the only thing here that is.
func (l *Looper) Wake() { C.ALooper_wake(l.ptr) }
