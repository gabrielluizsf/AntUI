//go:build android

package app

/*
#include <stdlib.h>
#include "glue.h"
*/
import "C"

import (
	"unsafe"
)

// antuiOnSaveInstanceState is called by the framework, on the UI thread,
// when it wants the app's state written down.
//
// The block handed back is freed by the framework with free, so it has to
// come from malloc — a Go slice's memory would be freed by something that
// does not own it, which is a corruption rather than a crash.
//
//export antuiOnSaveInstanceState
func antuiOnSaveInstanceState(_ *C.ANativeActivity, size *C.size_t) unsafe.Pointer {
	*size = 0
	saveMu.Lock()
	f := saveFn
	saveMu.Unlock()
	if f == nil {
		return nil
	}
	b := f()
	if len(b) == 0 {
		return nil
	}
	p := C.malloc(C.size_t(len(b)))
	if p == nil {
		return nil
	}
	copy(unsafe.Slice((*byte)(p), len(b)), b)
	*size = C.size_t(len(b))
	return p
}
