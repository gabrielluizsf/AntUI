//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

import "errors"

// attach finds this thread's connection, joining it to the machine if it is
// not already, and reports whether it has to be detached again.
func attach(m *C.JavaVM) (env *C.JNIEnv, detach bool, err error) {
	switch C.mw_GetEnv(m, &env, C.JNI_VERSION_1_6) {
	case C.JNI_OK:
		return env, false, nil
	case C.JNI_EDETACHED:
		if C.mw_AttachCurrentThread(m, &env) != C.JNI_OK {
			return nil, false, errors.New("jni: cannot attach this thread to the virtual machine")
		}
		return env, true, nil
	}
	return nil, false, errors.New("jni: the virtual machine does not speak JNI 1.6")
}
