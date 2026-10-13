//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "errors"

// findAppLoader reaches the activity's class loader and brings it back as a
// global reference, with loadClass on it.
func (e *Env) findAppLoader() (C.jobject, C.jmethodID, error) {
	act := Activity()
	if act.IsNil() {
		return 0, nil, ErrNoVM
	}
	// The activity's class, then that class's own class, which is
	// java/lang/Class — reached this way rather than by name so that it
	// works on a thread FindClass would fail on.
	actClass := C.mw_GetObjectClass(e.ptr, act.ref)
	if actClass == 0 {
		return 0, nil, errors.New("jni: the activity has no class")
	}
	classClass := C.mw_GetObjectClass(e.ptr, C.jobject(actClass))
	get := e.rawMethod(classClass, "getClassLoader", "()Ljava/lang/ClassLoader;")
	if get == nil {
		return 0, nil, errors.New("jni: java.lang.Class has no getClassLoader")
	}
	ldr := C.mw_CallObject(e.ptr, C.jobject(actClass), get, nil)
	if err := e.check(); err != nil {
		return 0, nil, err
	}
	if ldr == 0 {
		return 0, nil, errors.New("jni: the activity's class has no loader")
	}
	ldrClass := C.mw_GetObjectClass(e.ptr, ldr)
	load := e.rawMethod(ldrClass, "loadClass", "(Ljava/lang/String;)Ljava/lang/Class;")
	if load == nil {
		return 0, nil, errors.New("jni: the class loader has no loadClass")
	}
	global := C.mw_NewGlobalRef(e.ptr, ldr)
	if global == 0 {
		return 0, nil, errors.New("jni: cannot hold on to the class loader")
	}
	return global, load, nil
}
