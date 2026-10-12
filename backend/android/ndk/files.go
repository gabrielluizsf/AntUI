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

// Files lists the assets directly inside a directory.
//
// **It lists files and not directories**, which is a limitation of this API
// and not of the package: there is no NDK call that reports a subdirectory.
// The assets package above this one fills that in through Java, where the
// same list does include them.
func (m *AssetManager) Files(dir string) ([]string, error) {
	if m == nil || m.ptr == nil {
		return nil, errors.New("ndk: there is no asset manager")
	}
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	d := C.AAssetManager_openDir(m.ptr, cdir)
	if d == nil {
		return nil, ErrNoAsset
	}
	defer C.AAssetDir_close(d)

	var out []string
	for {
		name := C.AAssetDir_getNextFileName(d)
		if name == nil {
			break
		}
		out = append(out, C.GoString(name))
	}
	return out, nil
}
