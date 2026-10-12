//go:build android

package app

/*
#include <android/looper.h>
#include "glue.h"
*/
import "C"

import (
	"syscall"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// startUIRunner sets the pipe up. It runs on the UI thread, inside
// ANativeActivity_onCreate, because that is the only place this thread's
// looper can be got at — a looper belongs to its thread and there is no way
// to ask for another one's.
func startUIRunner() {
	uiLooper = C.ALooper_forThread()
	if uiLooper == nil {
		ndk.Errorf("the UI thread has no looper; nothing can be run on it")
		return
	}
	var fds [2]int
	// Non-blocking, so draining the pipe cannot stall the UI thread, and
	// close-on-exec because a leaked descriptor outlives its usefulness.
	if err := syscall.Pipe2(fds[:], syscall.O_NONBLOCK|syscall.O_CLOEXEC); err != nil {
		ndk.Errorf("cannot make the UI work pipe: %v", err)
		return
	}
	if C.antui_watch_ui_pipe(uiLooper, C.int(fds[0])) != 1 {
		ndk.Errorf("the UI thread's looper would not watch the work pipe")
		syscall.Close(fds[0])
		syscall.Close(fds[1])
		return
	}
	uiMu.Lock()
	uiWrite = fds[1]
	uiMu.Unlock()
}
