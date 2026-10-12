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
	"unsafe"
)

// Open finds an asset by its path inside the package, with no leading slash:
// "levels/one.json".
func (m *AssetManager) Open(name string, mode int) (*Asset, error) {
	if m == nil || m.ptr == nil {
		return nil, errors.New("ndk: there is no asset manager")
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	a := C.AAssetManager_open(m.ptr, cname, C.int(mode))
	if a == nil {
		return nil, ErrNoAsset
	}
	return &Asset{ptr: a, size: int64(C.AAsset_getLength64(a))}, nil
}
