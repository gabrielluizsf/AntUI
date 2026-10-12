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
)

// Seek moves the read position, as an io.Seeker does.
func (a *Asset) Seek(offset int64, whence int) (int64, error) {
	if a.ptr == nil {
		return 0, errors.New("ndk: the asset is closed")
	}
	at := int64(C.AAsset_seek64(a.ptr, C.off64_t(offset), C.int(whence)))
	if at < 0 {
		return 0, errors.New("ndk: cannot seek in a compressed asset")
	}
	return at, nil
}
