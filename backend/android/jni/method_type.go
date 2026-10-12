//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Method identifies a method of a class. It is not a reference and does not
// need freeing, but it is only valid while its class is loaded — which, for
// a cached class, is forever.
type Method struct{ id C.jmethodID }
