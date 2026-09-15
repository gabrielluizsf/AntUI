//go:build android

package app

/*
#include <android/looper.h>
#include "glue.h"
*/
import "C"

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Most of Android's Java API may only be touched from the thread the
// activity was made on. Reading the window insets from anywhere else throws;
// changing the system bars from anywhere else is undefined.
//
// The usual answer is Activity.runOnUiThread, which takes a Runnable — a
// Java object, which needs a Java class, which needs a dex file in the APK.
// None of that is necessary. The UI thread already has a looper, that looper
// already knows how to wait on a file descriptor, and a pipe is a file
// descriptor. So the runner is a pipe and a callback, with no Java at all.
var (
	uiMu     sync.Mutex
	uiQueue  []func()
	uiWrite  int = -1
	uiLooper *C.ALooper
)

// ErrNoUIThread is what the runner reports before the activity exists.
var ErrNoUIThread = errors.New("app: there is no UI thread yet")

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

// OnUIThread reports whether the caller is already on the UI thread. It is
// what stops [RunOnUISync] deadlocking on itself.
func OnUIThread() bool {
	return uiLooper != nil && C.ALooper_forThread() == uiLooper
}

// RunOnUI puts f on the UI thread's queue and returns without waiting.
func RunOnUI(f func()) error {
	uiMu.Lock()
	if uiWrite < 0 {
		uiMu.Unlock()
		return ErrNoUIThread
	}
	uiQueue = append(uiQueue, f)
	fd := uiWrite
	uiMu.Unlock()
	// One byte is a nudge, not a message: the queue is what carries the
	// work, and a pipe that already has bytes in it wakes the looper just
	// the same.
	if _, err := syscall.Write(fd, []byte{0}); err != nil && err != syscall.EAGAIN {
		return err
	}
	return nil
}

// RunOnUISync runs f on the UI thread and waits for it to finish.
//
// Called from the UI thread it simply calls f, which is what keeps it from
// waiting for itself. From anywhere else it blocks — so it must not be
// called from anything the UI thread is itself waiting on, which in this
// library means it must not be called from inside an activity callback.
func RunOnUISync(f func()) error {
	if OnUIThread() {
		f()
		return nil
	}
	done := make(chan struct{})
	if err := RunOnUI(func() {
		defer close(done)
		f()
	}); err != nil {
		return err
	}
	<-done
	return nil
}

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
