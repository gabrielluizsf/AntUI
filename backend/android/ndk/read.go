//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

// Read fills p, as an io.Reader does.
func (a *Asset) Read(p []byte) (int, error) {
	if a.ptr == nil {
		return 0, errors.New("ndk: the asset is closed")
	}
	if len(p) == 0 {
		return 0, nil
	}
	n := int(C.AAsset_read(a.ptr, unsafe.Pointer(&p[0]), C.size_t(len(p))))
	switch {
	case n > 0:
		return n, nil
	case n == 0:
		return 0, io.EOF
	}
	return 0, errors.New("ndk: reading the asset failed")
}
