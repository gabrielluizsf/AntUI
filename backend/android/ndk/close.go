//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

// Close gives the asset back. A mapping handed out by [Asset.Bytes] stops
// being valid here.
func (a *Asset) Close() error {
	if a.ptr != nil {
		C.AAsset_close(a.ptr)
		a.ptr = nil
	}
	return nil
}
