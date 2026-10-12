//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// InputQueue is the stream of key presses and touches the system has for
// this window. It only produces anything while it is attached to a looper,
// and only on the thread that looper belongs to.
type InputQueue struct {
	ptr *C.AInputQueue
}
