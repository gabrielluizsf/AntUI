//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "unsafe"

// Native is one native method to attach to a Java class.
//
// Registering by hand is rarely needed. A Go function exported as
// Java_pkg_Class_method is found by the machine on its own, by name, the
// first time the method is called — which is how the shim in
// antui/backend/android/app works and is one fewer thing to get wrong. This is for
// the cases where the name cannot be used: a class the build renamed, or a
// method whose Go implementation is chosen at run time.
type Native struct {
	// Name is the method's name in Java.
	Name string
	// Sig is its signature, from [Sig].
	Sig string
	// Fn is the C address of the implementation. It is not an ordinary Go
	// function: it has to be a symbol cgo exported, taken as C.theSymbol
	// from the package that exported it.
	Fn unsafe.Pointer
}
