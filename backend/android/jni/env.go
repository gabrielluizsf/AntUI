//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

// Env is one thread's connection to the virtual machine, inside one local
// reference frame. It is only valid for the duration of the [Do] that made
// it — keeping one and using it later is a crash, not an error.
type Env struct {
	ptr *C.JNIEnv
}
