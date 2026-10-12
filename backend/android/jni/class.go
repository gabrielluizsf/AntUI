//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Class is a Java class, held as a global reference and cached, so that
// looking one up twice costs nothing the second time.
type Class struct{ ref C.jclass }
