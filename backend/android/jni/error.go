//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

var errNoFrame = errors.New("jni: cannot open a local reference frame")

// Error is a Java exception, carried across as a Go error.
type Error struct {
	// Class is the exception's class, dotted: "java.lang.SecurityException".
	Class string
	// Message is what getMessage said, which is often empty.
	Message string
}

// Error is the Java exception as a Go error: the class it was, and what it
// said.
func (e *Error) Error() string {
	if e.Message == "" {
		return "java: " + e.Class
	}
	return "java: " + e.Class + ": " + e.Message
}

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

// describe calls getClass on the throwable and then getName on the result.
// It uses the raw calls rather than the checked ones, because this is what
// runs when something has already gone wrong and a failure here must not
// turn into another error.
func (e *Env) describe(t C.jobject, get, sig, then string) string {
	cls := C.mw_GetObjectClass(e.ptr, t)
	if cls == 0 {
		clearPending(e.ptr)
		return ""
	}
	m := e.rawMethod(cls, get, sig)
	if m == nil {
		return ""
	}
	obj := C.mw_CallObject(e.ptr, t, m, nil)
	clearPending(e.ptr)
	if obj == 0 {
		return ""
	}
	return e.text(obj, then, "()Ljava/lang/String;")
}

// text calls a no-argument method that returns a String, and reads it. It
// swallows anything that goes wrong, for the same reason describe does.
func (e *Env) text(o C.jobject, name, sig string) string {
	cls := C.mw_GetObjectClass(e.ptr, o)
	if cls == 0 {
		clearPending(e.ptr)
		return ""
	}
	m := e.rawMethod(cls, name, sig)
	if m == nil {
		return ""
	}
	s := C.mw_CallObject(e.ptr, o, m, nil)
	clearPending(e.ptr)
	if s == 0 {
		return ""
	}
	out, err := e.GoString(Object{ref: s})
	if err != nil {
		return ""
	}
	return out
}

func (e *Env) rawMethod(cls C.jclass, name, sig string) C.jmethodID {
	cname, csig := C.CString(name), C.CString(sig)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(csig))
	m := C.mw_GetMethodID(e.ptr, cls, cname, csig)
	if m == nil {
		clearPending(e.ptr)
	}
	return m
}
