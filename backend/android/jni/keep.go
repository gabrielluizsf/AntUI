//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"runtime"
)

// keep stops the collector taking the argument array away while the call
// that is reading it is still running.
func keep(a []C.jvalue) { runtime.KeepAlive(a) }
