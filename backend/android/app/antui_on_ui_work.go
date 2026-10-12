//go:build android

package app

/*
#include <android/looper.h>
#include "glue.h"
*/
import "C"

import (
	"syscall"
	"unsafe"
)

// antuiOnUIWork runs on the UI thread whenever the pipe has something in it.
// Returning 1 keeps the descriptor registered; returning 0 would take it off
// the looper and the runner would silently stop working.
//
//export antuiOnUIWork
func antuiOnUIWork(fd C.int, events C.int, data unsafe.Pointer) C.int {
	var drain [64]byte
	for {
		n, err := syscall.Read(int(fd), drain[:])
		if n < len(drain) || err != nil {
			break
		}
	}
	for {
		uiMu.Lock()
		if len(uiQueue) == 0 {
			uiMu.Unlock()
			return 1
		}
		next := uiQueue[0]
		uiQueue = uiQueue[1:]
		uiMu.Unlock()
		next()
	}
}
