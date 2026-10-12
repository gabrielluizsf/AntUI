//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

import (
	"runtime"
	"sync"
)

// Hold attaches the calling goroutine's thread and keeps it attached until
// the returned function is called. [Do] inside the hold is then only a
// lookup rather than an attach, which is what makes it cheap enough to use
// in a frame loop.
//
// It is also how a **sequence** of calls is kept on one thread, which some
// of Java's API quietly requires: a SQLite transaction belongs to the thread
// that began it, and a lock taken on one thread cannot be released on
// another. Go moves a goroutine between threads whenever it likes, and each
// [Do] only pins for its own duration — so anything that must happen on one
// thread from beginning to end has to hold.
//
// The goroutine stays pinned to its thread for the whole hold. **The
// returned function must be called before the goroutine ends**: a thread
// that dies while still attached takes the whole process down with it, and
// the message the runtime prints does not say so.
func Hold() (release func(), err error) {
	m := machine()
	if m == nil {
		return nil, ErrNoVM
	}
	runtime.LockOSThread()
	_, detach, err := attach(m)
	if err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			if detach {
				C.mw_DetachCurrentThread(m)
			}
			runtime.UnlockOSThread()
		})
	}, nil
}
