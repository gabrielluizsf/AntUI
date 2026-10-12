//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// check turns a pending exception into an error, and clears it.
//
// It is called after **every** call into Java, and that is not caution: JNI
// does not stop when a method throws. The exception stays pending, the next
// call into the machine fails or aborts, and the message names whatever was
// being done then rather than what actually went wrong. Clearing it here is
// what makes an error mean what it says.
func (e *Env) check() error {
	if C.mw_ExceptionCheck(e.ptr) == C.JNI_FALSE {
		return nil
	}
	t := C.mw_ExceptionOccurred(e.ptr)
	// Clear first. While an exception is pending almost nothing may be
	// called, including the calls needed to find out what it was.
	C.mw_ExceptionClear(e.ptr)
	if t == 0 {
		return &Error{Class: "java.lang.Throwable"}
	}
	err := &Error{
		Class:   e.describe(C.jobject(t), "getClass", "()Ljava/lang/Class;", "getName"),
		Message: e.text(C.jobject(t), "getMessage", "()Ljava/lang/String;"),
	}
	if err.Class == "" {
		err.Class = "java.lang.Throwable"
	}
	C.mw_DeleteLocalRef(e.ptr, C.jobject(t))
	return err
}
