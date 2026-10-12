//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import (
	"sync"
)

// There is exactly one activity in a process, so it is a package variable
// rather than something threaded through every callback. It is written on
// the UI thread in onCreate and onDestroy and read there and on the app's
// goroutine, so it is behind a lock.
var (
	stateMu  sync.Mutex
	theApp   *App
	activity *C.ANativeActivity
)
