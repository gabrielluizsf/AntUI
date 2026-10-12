//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

func machine() *C.JavaVM {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return vm
}
