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

func formatInt(f *C.AMediaFormat, key string) (int, bool) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var out C.int32_t
	if !C.AMediaFormat_getInt32(f, k, &out) {
		return 0, false
	}
	return int(out), true
}
