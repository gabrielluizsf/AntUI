//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Object is the class itself as an object — a java.lang.Class instance,
// which is what a constructor like Intent(Context, Class) takes. A class
// reference and an object reference are the same thing to the machine; only
// Go's types tell them apart.
func (c Class) Object() Object { return Object{ref: C.jobject(c.ref)} }
