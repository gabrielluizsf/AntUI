//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

import "sync"

// The virtual machine and the activity, set once by the app package when
// Android creates the activity. They are the only globals: everything else
// is per-thread or per-call by nature.
var (
	stateMu  sync.RWMutex
	vm       *C.JavaVM
	activity C.jobject // a global reference, so it outlives onCreate
)
