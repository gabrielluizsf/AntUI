//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
)

// Class finds a class by its name with slashes: "android/os/Build",
// "java/lang/String".
//
// It looks twice, and the second look is the point of this function.
// FindClass uses the class loader of whatever native code is on the stack —
// and on a thread this library attached, there is none, so it falls back to
// the *system* loader, which can see the framework and cannot see a single
// class from the app. The symptom is a ClassNotFoundException for a class
// that is plainly in the APK, on some threads and not others. The fix is to
// ask the app's own loader, which is what happens when the first look fails.
func (e *Env) Class(name string) (Class, error) {
	classMu.RLock()
	c, ok := classes[name]
	classMu.RUnlock()
	if ok {
		return Class{ref: c}, nil
	}

	ref := e.findClass(name)
	if ref == 0 {
		var err error
		if ref, err = e.loadClass(name); err != nil {
			return Class{}, err
		}
	}
	global := C.jclass(C.mw_NewGlobalRef(e.ptr, C.jobject(ref)))
	if global == 0 {
		return Class{}, errors.New("jni: cannot hold on to class " + name)
	}
	classMu.Lock()
	// Another goroutine may have got there first, in which case its
	// reference is the one everybody else already has.
	if was, ok := classes[name]; ok {
		classMu.Unlock()
		C.mw_DeleteGlobalRef(e.ptr, C.jobject(global))
		return Class{ref: was}, nil
	}
	classes[name] = global
	classMu.Unlock()
	return Class{ref: global}, nil
}
