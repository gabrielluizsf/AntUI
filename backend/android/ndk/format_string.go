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

func formatString(f *C.AMediaFormat, key string) string {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var out *C.char
	if !C.AMediaFormat_getString(f, k, &out) {
		return ""
	}
	return C.GoString(out)
}
