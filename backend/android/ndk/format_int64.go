//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import (
	"unsafe"
)

func formatInt64(f *C.AMediaFormat, key string) (int64, bool) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var out C.int64_t
	if !C.AMediaFormat_getInt64(f, k, &out) {
		return 0, false
	}
	return int64(out), true
}
